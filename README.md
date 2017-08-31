# mailenc
Encryption of emails to selected recipients directly on the outgoing MTA. Encryption is to be done with S/MIME.

**Table of Contents**
- [Requirements](#requirements)
- [Usage](#usage)
- [Command line flags](#command-line-flags)
- [Example with Postfix](#example-with-postfix)


# Requirements

mailenc only works with linux systems and openssl installed.
openssl currently is necessary cause I was too lazy to implement it via libopenssl in cgo.
Maybe I'll fix it at some time.


# Usage

You can can use mailenc to pipe mails from e.g. postfix into. It will encrypt the mails with the smime certificate of the mail recipient. Certificates have to be files in the form {user}@{domain.tld}.crt
The directory where the certificates are stored can be defined.


# Command line flags

```sh
$ go run main.go -h
```
should output
```
Usage of mailenc:
  -addr string
    	smtp address to send to (default "localhost:10025")
  -dir string
    	smime cert directory (default "/etc/postfix/smime")
  -rcpt string
    	mail recipient (default "recipient@example.com")
  -sender string
    	mail sender (default "sender@example.com")
```


# Example with Postfix

**/etc/postfix/master.cf**

Creating our milter to send mails through mailenc:

```
encryptsmime unix - n n - 2 pipe
    flags=Rq user=filter null_sender=
    argv=/usr/local/bin/mailenc -from ${sender} -to ${recipient} -dir /etc/postfix/smime -addr localhost:10025
```

Next create a new smtp listener so mails won't loop indefinitely. This is where mailenc will send mails to.

```
localhost:10025      inet  n       -       n       -       -       smtpd
       -o receive_override_options=no_milters
       -o smtpd_recipient_restrictions=permit_mynetworks,reject

```

**/etc/postfix/main.cf**


mailenc requires all recipients individually.

```
encryptsmime_destination_recipient_limit = 1
```

Add the milter for specific users.
Be careful though we're overriding the default behaviour of smtpd_recipient_restrictions!


```
smtpd_recipient_restrictions =  check_recipient_access hash:/etc/postfix/smime_access

```

Create the file **/etc/postfix/smime_access**
```
user@example.org	FILTER	encryptsmime:
```

and execute `postmap /etc/postfix/smime_access`



You can also do this without applying the milter to individual users, but this may hurt your performance if your system has to send hundreds of mails per second, cause every mail has to pass the milter and therefor is queued twice.
Just add `-o content_filter=encryptsmime:dummy` in /etc/postfix/master.cf to this line `smtp inet n - - - - smtpd`:

```
smtp inet n - - - - smtpd -o content_filter=encryptsmime:dummy
```
