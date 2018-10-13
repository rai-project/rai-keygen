# RAI KeyGen [![Build Status](https://travis-ci.org/rai-project/rai-keygen.svg?branch=master)](https://travis-ci.org/rai-project/rai-keygen)

## Generate Key

Generates profiles to be used with the rai client

### Synopsis

Generates a profile file that needs to be placed in `~/.rai_profile` (linux/OSX) or `%HOME%/.rai_profile` (Windows -- for me this is `C:\Users\abduld\.rai_profile`). The rai client reads these configuration files to authenticate the user. A seed (specified by `secret`) is used to generate secure credentials

```
rai-keygen
```

### Options

```
  -c, --color              Toggle color output.
  -d, --debug              Toggle debug mode.
  -e, --email string       The email to generate the key for.
  -r, --role string        The role (or privaleges) of the user (student, power, etc...).
  -f, --firstname string   The firstname to generate the key for.
  -l, --lastname string    The lastname to generate the key for.
  -s, --secret string      The application secret key.
  -u, --username string    The username to generate the key for.
  -v, --verbose            Toggle verbose mode.
```

## Email Keys

Creates keys for each person in a CSV file and emails the key to them.

### Synopsis

Creates keys for each user in the students list and emails it to the them. The student list must be formated as a CSV file and be of the form `lastname, firstname, username, email, role, affiliation`. The config file must be configured with the propper mailing credentials.

```
rai-keygen emailkeys
```

### Options

```
  -l, --studentlist string   The student list is the file that contains a list of all the students in csv format [lastname,filename,email].
      --subject string       The subject line for the email sent.
  -t, --template string      The email template file to use when sending emails to the students.
```

### Options inherited from parent commands

```
  -c, --color           Toggle color output.
  -d, --debug           Toggle debug mode.
  -s, --secret string   The application secret key.
  -v, --verbose         Toggle verbose mode.
```

## License

NCSA/UIUC © [Abdul Dakkak](http://impact.crhc.illinois.edu/Content_Page.aspx?student_pg=Default-dakkak)

[github issue manager]: https://github.com/rai-project/rai/issues
