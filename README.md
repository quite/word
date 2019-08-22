
For looking up words, I like to keep an
[RFC2229](https://tools.ietf.org/html/rfc2229) dictd answering on (only)
`localhost` with a few of databases installed. This lookup tool is presently
tailored to myself. It does automagic paging!

By default, word talks to `localhost:2628` and looks up the word, its first
argument, in these locally installed databases: gcide wn foldoc moby-thesaurus.

On Arch Linux, I install and run dictd like so (side-stepping the sysv-style
/etc/conf.d/dictd configuration that dictd's systemd service file still uses):

    pacman -S dictd dict-gcide dict-wn dict-foldoc moby-thesaurus
    cp --parents dictd.service.d/override.conf /etc/systemd/system/
    systemctl enable --now dictd

Stay local.
