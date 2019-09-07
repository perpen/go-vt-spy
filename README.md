Trivial terminal session recorder

It only records stdout (which ofc gets stderr as well in a tty), logging
stdin as well would be trivial. 

Start with `./vt-spy 2> session.log`
