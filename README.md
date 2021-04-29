# web-remote-control-bedroom
A web GUI to control my Raspberry Pi Blinkt over the Internet!

# Overview
This is a very (very!) basic interface to configure the LEDs in a Blinkt piHat installed in my room.

The ultimate goal of this project is to experiment with serverless functions (aka Lambda functions), MQ and adressable LEDs. This project is more of a study than actually a product - feel free to adapt it, extend it and improve, but, keep in mind that this project was not made having performance in mind.

# What does this do?
This project is a website that has 8 circles (representing the 8 LEDs in BLinkt) that can be customized for the color and brightness of you choice. For that, just select the color in the color select button and place the slider at the desired brightness. Once they are set, click on the circle (LED) which you want to have those properties. Black circles mean the LED will be turn off. Once all circles are configured as you want, click the `Update Blinkt` button and the Blinkt piHat will set the LEDs accordingly.

This is a push-only website, meaning that if Blinkt already has been configured before loading the webpage, it will not reflect the current state of the LEDs in Blinkt.

# How to use it
This is a Netlify-ready project, meaning you can use it as is and have a working website without any customization. The MQ is the public Hive MQ, so do not use this project as is, unless for non-production, insecure projects and with no sensible data at all.

To deploy this project in Netlify, just follow their instructions at: [Netlify CLI doc](https://docs.netlify.com/cli/get-started/).

Of course, Blinkt needs to have the correct program to read the entries in MQ and set the LEDs correctly. For that, you need another project: [blinkt-controlled-mq](https://github.com/otaviokr/blinkt-controlled-mq.git).

# If you want a simpler alternative
This project is an evolution of another, that would provide the webpage running locally in the Raspberry Pi where Blinkt piHat is installed. Check it out at: [blinkt-web-ui](https://github.com/otaviokr/blinkt-web-ui).
