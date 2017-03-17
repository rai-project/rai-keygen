package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/auth"
	"github.com/rai-project/auth/provider"
	"github.com/rai-project/email/mailgun"
	"github.com/spf13/cobra"
)

var (
	studentListFileName   string
	emailTemplateFileName string
	emailSubjectLine      string
)

var emailKeysCmd = &cobra.Command{
	Use:     "emailkeys",
	Aliases: []string{"email"},
	Short:   "Creates keys for each person in a CSV file and emails the key to them.",
	Long: "Creates keys for each user in the students list and emails it to the them. " +
		"The student list must be formated as a CSV file and be of the form `lastname, firstname, username, email, affiliation`. " +
		"The config file must be configured with the propper mailing credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := log.WithField("cmd", "emailkeys")

		if studentListFileName == "" || !com.IsFile(studentListFileName) {
			return errors.New("The student file list has not been found.")
		}
		if emailTemplateFileName != "" && !com.IsFile(emailTemplateFileName) {
			return errors.Errorf("cannot find the email template file %v", emailTemplateFileName)
		}

		studentFile, err := os.Open(studentListFileName)
		if err != nil {
			return errors.Wrap(err, "Failed to open the student list file")
		}
		defer studentFile.Close()

		var emailTemplateFileBytes []byte

		if emailTemplateFileName == "" {
			emailTemplateFileBytes, err = _escFSByte(false, "emailkey.template")
		} else {
			emailTemplateFileBytes, err = ioutil.ReadFile(emailTemplateFileName)
		}
		if err != nil {
			return errors.Wrap(err, "Failed to read the email template file")
		}

		emailTemplateFileContent := string(emailTemplateFileBytes)

		emailTemplate, err := template.New("email_template").Parse(emailTemplateFileContent)

		mail, err := mailgun.New()
		if err != nil {
			return err
		}

		r := csv.NewReader(studentFile)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Error(%v):: cannot read record in the student file list.\n", err)
				continue
			}
			if len(record) != 5 {
				fmt.Printf("Error(%v):: cannot read record in the student file list. "+
					"The format must be [lastname, firstname, username, email, affiliation]\n", err)
				continue
			}

			prof, err := provider.New(
				auth.Lastname(record[0]),
				auth.Firstname(record[1]),
				auth.Username(record[2]),
				auth.Email(record[3]),
				auth.Affiliation(record[4]),
			)
			if err != nil {
				return err
			}

			user := prof.Info()
			profFileContent, err := yaml.Marshal(user)
			if err != nil {
				fmt.Printf("Error(%v):: cannot create profile for user.\n", err)
				continue
			}

			templateParams := struct {
				auth.ProfileBase
				ProfileContent     string
				ProfileFileContent string
			}{
				ProfileBase:        user,
				ProfileContent:     user.String(),
				ProfileFileContent: string(profFileContent),
			}

			emailBody := new(bytes.Buffer)
			err = emailTemplate.Execute(emailBody, templateParams)
			if err != nil {
				fmt.Printf("Error(%v):: failed to create email message body. \n", err)
				continue
			}
			err = mail.Send(templateParams.Email, emailSubjectLine, emailBody.String())
			if err != nil {
				fmt.Printf("Failed to send email to %s.\n", user.Email)
				continue
			}

			log.WithField("first_name", user.Firstname).
				WithField("last_name", user.Lastname).
				WithField("user_name", user.Username).
				WithField("secret", user.SecretKey).
				WithField("access", user.AccessKey).
				Infof("Email was successfully sent to %s.\n", user.Email)
		}
		return nil
	},
}

func init() {
	emailKeysCmd.PersistentFlags().StringVarP(&studentListFileName, "studentlist", "l", studentListFileName,
		"The student list is the file that contains a list of all the students in csv format [lastname,filename,email].")
	emailKeysCmd.PersistentFlags().StringVarP(&emailTemplateFileName, "template", "t", "",
		"The email template file to use when sending emails to the students.")
	emailKeysCmd.PersistentFlags().StringVar(&emailSubjectLine, "subject", emailSubjectLine,
		"The subject line for the email sent.")

	RootCmd.AddCommand(emailKeysCmd)
}
