# pap :microphone: :speaker:

[![Build Status](https://travis-ci.org/tomyl/pap.svg?branch=master)](https://travis-ci.org/tomyl/pap)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomyl/pap)](https://goreportcard.com/report/github.com/tomyl/pap)

A simple pulseaudio profile manager. Makes it easy to switch between pairs of sources/sinks. 

# Usage

Use `pap -next-sink` to cycle between sinks. The chosen sink will be marked as
default and active playback streams will switch to it. If the card providing the sink
has a source, that will become the default source. Use `pap -list-sinks` to
show available sinks.

```bash
$ pap -next-sink
Activated NuForce µDAC 2 Analog Stereo.
$ pap -list-sinks                         
Built-in Audio Digital Stereo (HDMI)
Built-in Audio Analog Stereo
ClearChat Pro USB Analog Stereo
NuForce µDAC 2 Analog Stereo [current]
```

You can also save custom source/sink pairs as a profile:

```bash
$ pavucontrol # choose default source and sink
$ pap -add headset   
Added profile headset.
$ pap -list
headset [current]
laptop
nuforce
$ pap -next
Activated laptop.
```

Tip: use e.g. `pap -next -notify` to show messages in desktop notifcations instead of standard output.

# Installation

```bash
$ apt-get install libpulse-dev
$ go get github.com/tomyl/pap
$ ~/go/bin/pap -help
```
