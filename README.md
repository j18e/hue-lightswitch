# Hue Lightswitch
Use cheap 433Mhz buttons to control Philips Hue lights

## Hardware
- 433Mhz lightswitches [like this](https://www.ebay.com/itm/273799355401?hash=item3fbfb69c09:g:qrYAAOSwo61crBAe)
- 433Mhz receiver [like this](https://www.ebay.com/itm/273512863339?hash=item3faea3166b:g:dGkAAOSw-TFd99PM)
- a computer to run the software with the receiver connected (eg raspberry pi
  zero)
- a functioning Philips Hue system with a bridge

## Setting up
Install the `rtl_433` software to both your server and your development machine.
Homebrew and Aptitude should both have it. If in doubt its repo can be found
[here](https://github.com/merbanan/rtl_433). Run `rtl_433` with the 433Mhz
receiver connected via USB and click the different buttons on your switches.
The button presses should produce terminal output in the program. Make note of
the switch name and data fields, and modify `config.yml` with their values.

Complete your modifications of `config.yml` by creating mappings that reflect
the actions you want to occur when you press the buttons. You can get a dump of
the light data on your Hue Bridge by using the `lights` argument when running
the `hue-lightswitch` program.

Get a user token from your Hue Bridge by using their HTTP API and pushing its
physical button. Put this token in your cloned repo in a file called `.token`.

Populate a file called `.env` in the repo directory with the variables listed in
the top of the `Makefile`.

Run `make all`.
