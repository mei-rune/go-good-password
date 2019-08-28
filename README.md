# Go password checker

A simple library to check how strong a password is.

## Do you need this?

This aims to strike a balance between very strict password policies and provide
a simple idea of how strong it is.

Please if you haven't read the recommendations from NIST read:
https://pages.nist.gov/800-63-3/sp800-63b.html#memsecret

They boil down to require 8 characters at minimum and it is recommended to block
disclosed passwords. This library will return a very low score for passwords
less than 8 characters (including unicode characters).

This is partly aiming to be lightweight -- if you're running a web service you
may wish to consider "zxcvbn" which covers many more things. (And ideally give
realtime feedback via JavaScript, my initial use case for this was a CLI
utility).

## Implementing

    import "github.com/dgl/go-good-password"

Then to use it:

    // Put common words here, both from user's name, email and your service's name.
    extra := good_password.ExtractWords(user, email, "your-service-name")

    score, info := good_password.Check(password, extra)

    if score < 1 {
      // Don't allow passwords with really common words or such.
      fmt.Printf("%v password, pick a better one (%v)\n", score, info)
      return
    }

    // Otherwise guilt the user into picking a better password, maybe, but let
    // them do whatever. (Also you could show them info, but be careful about
    // logging it.)
    fmt.Printf("%v password!\n", score)

## Details

Score is an integer, an 8 character password with lowercase letters will score 1
(aka "terrible"), unless it has common words, repeats, patterns or sequences.
Increasing the length and using multiple types of character will increase the
score. The `info` returned will explain what was positive or negative about the
password (see API docs).

To meet the NIST recommendations above simply block a score less than one. For
more strict password requirements you can require higher scores.

Examples of scores:

* 1 _("terrible")_: "something" (one type)
* 2 _("weak")_: "somethin1", "somethingnew" (two types)
* 3 _("okay")_: "Somethin1", "somethinglonger" (three types)
* 4 _("good")_: "Someth!n1", "somethingmuchlonger" (four types)
* >=5 ("strong"): "Someth!n10", "correct horse battery staple" (five types)

It's also possible to score more by having a longer password. This means
xkcd.com/936 passwords are allowed with a score of at least 4 (provided they are
16 characters or more).

I suggest not requiring a score of more than 4 except in very specific cases
(this allows the length of the password alone to be enough and therefore doesn't
impose arbitrary rules on the user).

### Unicode

This library correctly handles unicode for lengths of passwords. It does not
perform normalisation on the password. This is your responsibility, see
https://blog.golang.org/normalization.

It may be as simple as:

    storePassword := norm.NFKC.String(password)

## Other implementations

- https://github.com/nbutton23/zxcvbn-go -- for my case it was too big (about
  1MB added to binary).
