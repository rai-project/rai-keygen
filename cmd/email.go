package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/auth"
	"github.com/rai-project/email/mailgun"
	"github.com/spf13/cobra"
)

type User struct {
	LastName  string `toml:"-"`
	FirstName string `toml:"-"`
	Username  string `toml:"username"`
	Email     string `toml:"-"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
}

var (
	// studentListFileName = "/Users/abduld/Code/wbgo/utils/408users.csv"
	studentListFileName   = ""
	emailTemplateFileName string
	emailSubjectLine      = "ECE 408 Remote Development Resource Information"
)

var emailKeysCmd = &cobra.Command{
	Use:   "emailkeys",
	Short: "Creates keys for each user in the students list and emails it to the them.",
	Long: "Creates keys for each user in the students list and emails it to the them. " +
		"The student list must be formated as a CSV file and be of the form firstname,lastname,email ." +
		"Another parameter that's needed is the mailgun key.",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := log.WithField("cmd", "emailkeys")

		if studentListFileName == "" || !com.IsFile(studentListFileName) {
			return errors.New("The student file list has not been found.")
		}
		if emailTemplateFileName == "" {
			srcPath, err := sourcepath.AbsoluteDir()
			if err != nil {
				return errors.New("The email template file has not been set.")
			}
			emailTemplateFileName = filepath.Join(srcPath, "rai-emailkeys", "emailkey.template")
		}
		if !com.IsFile(emailTemplateFileName) {
			return errors.New("The email template file has not been found.")
		}

		studentFile, err := os.Open(studentListFileName)
		if err != nil {
			return errors.Wrap(err, "Failed to open the student list file")
		}
		defer studentFile.Close()

		emailTemplateFileBytes, err := ioutil.ReadFile(emailTemplateFileName)
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
			if len(record) != 4 {
				fmt.Printf("Error(%v):: cannot read record in the student file list. "+
					"The format must be [lastname, firstname, username, email]\n", err)
				continue
			}
			username := record[2]
			accessKey, secretKey, err := auth.Hash(username)
			if err != nil {
				return err
			}
			user := User{
				LastName:  record[0],
				FirstName: record[1],
				Username:  username,
				Email:     record[3],
				SecretKey: secretKey,
				AccessKey: accessKey,
			}

			emailBody := new(bytes.Buffer)
			err = emailTemplate.Execute(emailBody, user)
			if err != nil {
				fmt.Printf("Error(%v):: failed to create email message body. \n", err)
				continue
			}
			err = mail.Send(user.Email, emailSubjectLine, emailBody.String())
			if err != nil {
				fmt.Printf("Failed to send email to %s.\n", user.Email)
				continue
			}

			log.WithField("first_name", user.FirstName).
				WithField("last_name", user.LastName).
				WithField("user_name", user.Username).
				WithField("secret", user.SecretKey).
				WithField("access", user.AccessKey).
				Infof("Email was successfully sent to %s.\n", user.Email)
		}
		return nil
	},
}

func initEmailKeys() {
	emailKeysCmd.PersistentFlags().StringVarP(&studentListFileName, "studentlist", "s", studentListFileName,
		"The student list is the file that contains a list of all the students in csv format [lastname,filename,email].")
	emailKeysCmd.PersistentFlags().StringVarP(&emailTemplateFileName, "template", "t", "",
		"The email template file to use when sending emails to the students.")
	emailKeysCmd.PersistentFlags().StringVar(&emailSubjectLine, "emailsubject", emailSubjectLine,
		"The subjectline for the email sent.")
}
