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

Behind the scenes, the API utilizes <a href="https://pkg.go.dev/github.com/gocolly/colly">Colly</a>, a Go package meant for web scraping, to parse up-to-date information every request. To scrape a webpage, Colly has a struct called a "collector". In order to speed up requests, once the API logs in to Home Access Center using certain credentials once, it will store this logged-in collector in memory, and utilize it to serve any future requests sent with the same credentials. The cache itself is implemented by the <a href="https://pkg.go.dev/github.com/jellydator/ttlcache/v3">TTLCache</a> package, a lightweight, yet fast cache which also allows for easy management of resources via expiration times. The API also utilizes Go's goroutines, which provide another boost in speed as it is trivial to parse HTML in parallel. The parsing is handled by <a href="https://pkg.go.dev/github.com/PuerkitoBio/goquery">GoQuery</a>, a package which allows for fast HTML parsing, and is compatible with Colly. The API depends on the <a href="https://pkg.go.dev/github.com/gofiber/fiber/v2">Fiber</a> web framework in order to handle requests and other parts of the API. It is bolstered by the fact that Fiber is optimized for speed. Finally, the API also makes use of <a href="https://pkg.go.dev/github.com/bytedance/sonic">Sonic</a> for blazingly fast JSON marshalling/unmarshalling. For generating documentation, the <a href="https://pkg.go.dev/github.com/swaggo/swag">Swag</a> package and Swagger are used. Further explanation of how the API gets the content from HAC can be found by examining the code in the _pkg_ and _app_ folders.

## Performance

## Credits

- The packages mentioned in the [How It Works](#how-it-works) section
- The [template](https://github.com/create-go-app/fiber-go-template) the package was based on
- The [Frisco ISD API](https://github.com/SumitNalavade/FriscoISDHACAPI) by Sumit Nalavade which gave inspiration on choosing Go and how to gather the information
- The [HACify](https://github.com/Threqt1/HACify) repo for providing the base logic for parsing the raw HTML
