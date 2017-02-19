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
	"github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/email"
	"github.com/spf13/cobra"
)

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
		slog := logrus.New()
		slog.Level = logrus.DebugLevel
		slog.Formatter = &logrus.TextFormatter{}
		log := slog.New().WithField("cmd", "emailkeys")

		if studentListFileName == "" || !com.IsFile(studentListFileName) {
			fmt.Println("Error:: the student file list has not been found.")
			return errors.New("The student file list has not been found.")
		}
		if emailTemplateFileName == "" {
			srcPath, err := sourcepath.AbsoluteDir()
			if err != nil {
				fmt.Println("Error:: the email template file has not been set.")
				return errors.New("The email template file has not been set.")
			}
			emailTemplateFileName = filepath.Join(srcPath, "rai-emailkeys", "emailkey.template")
		}
		if !com.IsFile(emailTemplateFileName) {
			fmt.Println("Error:: the email template file has not been found.")
			return errors.New("The email template file has not been found.")
		}

		studentFile, err := os.Open(studentListFileName)
		if err != nil {
			fmt.Printf("Error(%v):: failed to open the student list file.\n", err)
			return errors.Wrap(err, "Failed to open the student list file")
		}
		defer studentFile.Close()

		emailTemplateFileBytes, err := ioutil.ReadFile(emailTemplateFileName)
		if err != nil {
			fmt.Printf("Error(%v):: failed to read the email template file.\n", err)
			return errors.Wrap(err, "Failed to read the email template file")
		}
		emailTemplateFileContent := string(emailTemplateFileBytes)

		emailTemplate, err := template.New("email_template").Parse(emailTemplateFileContent)

		type User struct {
			LastName     string
			FirstName    string
			Email        string
			RAISecretKey string
			RAIAccessKey string
		}

		mail := email.New()

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
			if len(record) != 3 {
				fmt.Printf("Error(%v):: cannot read record in the student file list. "+
					"The format must be [lastname, firstname, email]\n", err)
				continue
			}
			key := keygen.New()
			user := User{
				LastName:     record[0],
				FirstName:    record[1],
				Email:        record[2],
				RAISecretKey: key.Secret,
				RAIAccessKey: key.Access,
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
				WithField("secret", user.RAISecretKey).
				WithField("access", user.RAIAccessKey).
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
