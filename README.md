# Chicago Gophers Meetup (July 17, 2018) on the Fn Project

## Resources

* Slides: [Google Slides link](https://docs.google.com/presentation/d/1HzdTt1xyr0fNRIZrhAE0ebLxEzXn9VmAhTbv8qreC6I/edit?usp=sharing)
* Project Site: https://fnproject.io
* Github: https://github.com/fnproject
* Go Tutorial: https://fnproject.io/tutorials/introduction
* Slack Channel: https://slack.fnproject.io

## Commands

### Install Fn Project
`curl -LSs https://raw.githubusercontent.com/fnproject/cli/master/install | sh`

### Start the Fn Project server
`fn start` 

### Start the Fn Project UI
`docker run --rm -it --link fnserver:api -p 4000:4000 -e "FN_API_URL=http://api:8080" fnproject/ui`

### Initialize a boiler plate Fn Project function (Go FDK)
`fn init --runtime go gofn`
`cd gofn`

### Test a function - CLI
`fn test`

### Run a function - CLI
`fn run` or `fn --verbose run`

### Run a function - CLI and pass input
`echo -n '{"name":"Bill"}' | fn run`

### Deploy a function to an "app"
`fn deploy --app gomeetup --local`

### View the calls for app "gomeetup"
`fn list calls gomeetup`

### View logs for a function call
`fn get logs gomeetup <call id>`


