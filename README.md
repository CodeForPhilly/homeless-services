# Resource Awareness Philly

Resource Awareness Philly is a open-source project that started at the Code for Philly SustyHack2015 from 10/17/15 - 10/18/15 in Philly. 

The goals of the project are

1. To create a comprehensive data set for all services for the homeless in Philly.
2. Provide web service REST apis to access and update the data set so that clients (web/mobile) can be developed.
3. Build a responsive web application that users in Philly can access on public library computers etc. as well as by services providers on desktop and mobile devices.
4. Iterative UI Design based on the needs and feedback of the homeless.
5. Increase awareness among volunteers to learn more about education & career opportunities related to coding & technology.

The software prototype includes

1. Web client app
2. REST APIs to access & update the data set.
3. Philly Homeless Resources data set as a Google doc.

This project is built with JavaScript/HTML/CSS on the front end and the Go programming language on the server side with DataStore for storage. All of it is served from or runs on Google App Engine.

A good intro to the Go programming language is here:
https://tour.golang.org/welcome/1

General documentation about Google App Engine's Go support is here:
https://cloud.google.com/appengine/docs/go/


## Set up for for this project

### Things to install

Install Go following these instructions
https://golang.org/doc/install

Install Git if you haven't already
https://git-scm.com/downloads

Install the Google App Engine SDK for Go
https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go

* git clone
* go get google.golang.org/appengine
* go get golang.org/x/net/context

### Running the application

Inside the rap directory, copy app.yaml.example to app.yaml

DO NOT check in your app.yaml because it can contain API and other private keys that shouldn't be shared.

Run this command where your app.yaml resides to start App Engine locally -> `goapp serve`

If you're playing with import, then this can be helpful to clear DataStore -> `goapp serve --clear_datastore`

From inside the rap directory you can run the tests -> `goapp test`

Deploy your code -> `goapp deploy`

The website will run locally on port 8080 and App Engine will run a local instance with access to things like DataStore and Memcache on 8000.

http://localhost:8080/

http://localhost:8000/

The file you'll want to import can be found in homeless-services/data/GeoCoded.csv

### If you have multiple versions of python

Things wil be harder.... But you can still get App Engine to work locally. This assumes you can normally run python 2.X with `python` from the command line. Move the app.yaml into the rap directory. Change the `basePath` in app.go to `./`. Then run this command from within the rap directory.

`python2 dev_appserver.py rap`

### The application may work locally but for it to work deployed, it'll need API Keys.

The Google Maps Key can created here; you'll probably need to do a little App Engine set up. The Geocoding key isn't used yet, but eventually should worked into the import process for locations that don't come with a Lat/Lon.
https://console.cloud.google.com/apis/credentials?project=yourappname

reCAPTCHA is used the feedback form. The code needed to validate it still needs to be written. This is a pretty development task if someone is interested.
https://www.google.com/recaptcha/admin#list

## (Optionally) Set up an IDE for Go

Download and install VS code -> https://code.visualstudio.com/Download

On Debian/Ubuntu, it can be installed with `sudo dpkg -i vscode-amd64.deb`

Install Go language support -> https://marketplace.visualstudio.com/items?itemName=lukehoban.Go

Get some go tools that will be used in the IDE (these are nice to have even with an IDE)

* go get -u github.com/rogpeppe/godef
* go get -u github.com/golang/lint/golint
* go get github.com/derekparker/delve/cmd/dlv

Open VS Code.

Click "Install Analysis Tools" in the lower right corner of VS Code.

If there's anything I missed, VS Code is pretty good about giving you little messages to fix things.

## License

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.

