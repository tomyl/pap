# pap

[![Build Status](https://travis-ci.org/tomyl/pap.svg?branch=master)](https://travis-ci.org/tomyl/pap)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomyl/pap)](https://goreportcard.com/report/github.com/tomyl/pap)

A simple pulseaudio profile manager. Usage:

```bash
$ pavucontrol # choose default source and sink
$ pap -add headset   
Added profile headset.
$ pap -list
headset*
laptop
nuforce
$ pap -next
Activated profile laptop.
$ pap -list
headset
laptop*
nuforce
```

Tip: use e.g. `pap -next -notify` to show messages in desktop notifcations instead of standard output.
