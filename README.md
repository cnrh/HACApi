Add shields.io badges, description: An API to serve homeaccess content | local setup | api docs | how it works | performance | tips | credits

<h1 align="center">HAC Information API</h1>

<h5 align="center">A fast and simple way to interact with <a href="https://www.powerschool.com/">PowerSchool Home Access Center</a></h5>

<div align="center">
  <a href=""><img src="https://img.shields.io/badge/Go-1.19.3-00ADD8?style=flat-square&logo=go" /></a>
</div>

## Description

The HAC Information API is an API written in Go meant to serve Home Access Center content quickly and in a format that's easy to work with.

## Local Setup

## API Docs

## How It Works

- Before the API is started, new documentation is generated using <a href="https://pkg.go.dev/github.com/swaggo/swag">Swag</a>, which parses comments in code to generate a Swagger template for the docs.

1. First, the API is started, middleware is registered, and routes are added. The framework the API uses <a href="https://pkg.go.dev/github.com/gofiber/fiber/v2">Fiber</a> to handle all of this behind the scenes.
2. Once a request is recieved at an endpoint, the body parameters are first validated. After that, the API will try to pull a logged-in <a href="https://pkg.go.dev/github.com/gocolly/colly">Colly</a> collector from the <a href="https://pkg.go.dev/github.com/jellydator/ttlcache/v3">TTLCache</a> for the provided credentials. If none are found, the API will try to use Colly to log into the provided HAC URL, and if successful, cache that collector for use in future requests.
3. The API will then use this logged-in collector to navigate to get the raw HTML for the requested data, once again using Colly. Once this raw HTML is recieved, the API uses <a href="https://pkg.go.dev/github.com/PuerkitoBio/goquery">GoQuery</a> to parse it, along with goroutines to parse in parallel for performance boosts.
4. Finally, the API will marshal this information into JSON using <a href="https://pkg.go.dev/github.com/bytedance/sonic">Sonic</a>, and send it back to the user.

- For more detailed insight into how the API queries and parses HAC information, consider reading the code in the _pkg_ and _app_ directories.

## Performance

## Credits

- The packages mentioned in the [How It Works](#how-it-works) section
- The [template](https://github.com/create-go-app/fiber-go-template) the package was based on
- The [Frisco ISD API](https://github.com/SumitNalavade/FriscoISDHACAPI) by Sumit Nalavade which gave inspiration on choosing Go and how to gather the information
- The [HACify](https://github.com/Threqt1/HACify) repo for providing the base logic for parsing the raw HTML
