---
date: 2017-03-17T00:48:33Z
title: "rai-keygen emailkeys"
slug: rai-keygen_emailkeys
url: /commands/rai-keygen_emailkeys/
---
## rai-keygen emailkeys

Creates keys for each person in a CSV file and emails the key to them.

### Synopsis


Creates keys for each user in the students list and emails it to the them. The student list must be formated as a CSV file and be of the form `lastname, firstname, username, email, affiliation`. The config file must be configured with the propper mailing credentials.

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

### SEE ALSO
* [rai-keygen](/commands/rai-keygen/)	 - Generates profiles to be used with the rai client

###### Auto generated by spf13/cobra on 17-Mar-2017
