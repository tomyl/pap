# pap :microphone: :speaker:

[![Build Status](https://travis-ci.org/tomyl/pap.svg?branch=master)](https://travis-ci.org/tomyl/pap)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomyl/pap)](https://goreportcard.com/report/github.com/tomyl/pap)

A simple pulseaudio profile manager. Makes it easy to switch between preconfigured pairs of sources/sinks. 

# Usage

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
$ pap -list -verbose
Loading /home/tomyl/.local/share/pap/profiles.json
headset
        source alsa_input.usb-Logitech_Logitech_USB_Headset-00.analog-mono (ClearChat Pro USB Analog Mono)
        sink   alsa_output.usb-Logitech_Logitech_USB_Headset-00.analog-stereo (ClearChat Pro USB Analog Stereo)
laptop*
        source alsa_input.pci-0000_00_1b.0.analog-stereo (Built-in Audio Analog Stereo)
        sink   alsa_output.pci-0000_00_1b.0.analog-stereo (Built-in Audio Analog Stereo)
nuforce
        source alsa_input.usb-NuForce_NuForce___DAC_2-01.analog-stereo (NuForce µDAC 2 Analog Stereo)
        sink   alsa_output.usb-NuForce_NuForce___DAC_2-01.analog-stereo (NuForce µDAC 2 Analog Stereo)
```

Tip: use e.g. `pap -next -notify` to show messages in desktop notifcations instead of standard output.

Alternatively, use the zero-configuration mode. `-list-auto` automatically builds one one profile per sink. If the card of the sink has any corresponding source, that becomes the profile source. Use `-next-auto` to switch between auto-generated profiles.

```bash
$ pap -list-auto                                  
Built-in Audio Digital Stereo (HDMI)
Built-in Audio Analog Stereo
ClearChat Pro USB Analog Stereo*
NuForce µDAC 2 Analog Stereo
$ pap -next-auto
Activated profile NuForce µDAC 2 Analog Stereo.
```

# Installation

```bash
$ apt-get install libpulse-dev
$ go get github.com/tomyl/pap
$ ~/go/bin/pap -help
```

# TODO

* [ ] When running `-next` and `-next-auto`, move all playback streams.
