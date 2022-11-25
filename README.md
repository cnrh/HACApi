<h1 align="center">HAC Information API</h1>

<h5 align="center">A fast and simple way to interact with <a href="https://www.powerschool.com/">PowerSchool Home Access Center</a></h5>

<div align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.19.3-00ADD8?style=flat-square&logo=go" /></a>
</div>

## Description

The HAC Information API is an API written in Go meant to serve Home Access Center content quickly and in a format that's easy to work with.
It covers the majority of the information HAC provides, including:

- Classwork (Per Marking Period)
- Interim Progress Reports (Per Date)
- Report Card(s)
- Transcript(s)
- Schedule(s)

With more features in the works, including:

- Week View (Per Day)
- Student Information
- Teacher Email Support
- Attendance
- Comment Legend
- `/multiple` endpoint for requesting multiple items at once

## Local Setup

1. Download [Go](https://go.dev/)
2. Clone the GitHub Repository
3. Rename the file `.env-example` to `.env` and fill in the required fields
4. Navigate into the folder, and run `go run main.go`

For Documentation:

1. Download [Swag](https://github.com/swaggo/swag) using `go install github.com/swaggo/swag/cmd/swag@latest`
2. Inside the main project, run

```bash
swag fmt
swag init
```

3. Once `swag init` finishes, the API automatically reserves the `/docs` endpoint for the docs. Navigate to it to view them.

## API Docs

Refer to the [API's Swagger Documentation](https://threqt1.github.io/HACApi/)

## How It Works

- Before the API is started, new documentation is generated using <a href="https://pkg.go.dev/github.com/swaggo/swag">Swag</a>, which parses comments in code to generate a Swagger template for the docs.

1. First, the API is started, middleware is registered, and routes are added. The framework the API uses <a href="https://pkg.go.dev/github.com/gofiber/fiber/v2">Fiber</a> to handle all of this behind the scenes.
2. Once a request is recieved at an endpoint, the body parameters are first validated. After that, the API will try to pull a logged-in <a href="https://pkg.go.dev/github.com/gocolly/colly">Colly</a> collector from the <a href="https://pkg.go.dev/github.com/jellydator/ttlcache/v3">TTLCache</a> for the provided credentials. If none are found, the API will try to use Colly to log into the provided HAC URL, and if successful, cache that collector for use in future requests.
3. The API will then use this logged-in collector to navigate to get the raw HTML for the requested data, once again using Colly. Once this raw HTML is recieved, the API uses <a href="https://pkg.go.dev/github.com/PuerkitoBio/goquery">GoQuery</a> to parse it, along with goroutines to parse in parallel for performance boosts.
4. Finally, the API will marshal this information into JSON using <a href="https://pkg.go.dev/github.com/bytedance/sonic">Sonic</a>, and send it back to the user.

- For more detailed insight into how the API queries and parses HAC information, consider reading the code in the _pkg_ and _app_ directories.

## Performance

### /login

| API Response Time (ms) | Login |
| ---------------------- | ----- |
|                        | 770   |
|                        | 715   |
|                        | 701   |
|                        | 822   |
|                        | 739   |
|                        | 681   |
|                        | 818   |
|                        | 761   |
|                        | 692   |
|                        | 743   |
| Average                | 744   |

### /classwork

Parameters:
(Not including username/password/base)

```json
{
  "markingPeriods": [1, 2, 3, 4, 5, 6]
}
```

| API Response Time (ms) | Without Login | With Login |
| ---------------------- | ------------- | ---------- |
|                        | 7530          | 7020       |
|                        | 8510          | 7180       |
|                        | 7050          | 7240       |
|                        | 7530          | 6090       |
|                        | 7930          | 6760       |
|                        | 7810          | 7600       |
|                        | 7650          | 7150       |
|                        | 7560          | 7360       |
|                        | 7640          | 6830       |
|                        | 7690          | 6780       |
| Average                | 7690          | 7001       |

### /ipr

| API Response Time (ms) | Without Login | With Login |
| ---------------------- | ------------- | ---------- |
|                        | 1270          | 582        |
|                        | 1264          | 562        |
|                        | 1230          | 559        |
|                        | 1501          | 586        |
|                        | 1464          | 572        |
|                        | 1355          | 569        |
|                        | 1349          | 761        |
|                        | 1258          | 589        |
|                        | 1327          | 560        |
|                        | 1214          | 562        |
| Average                | 1323          | 590        |

### /ipr/all

| API Response Time (ms) | Without Login | With Login |
| ---------------------- | ------------- | ---------- |
|                        | 1434          | 663        |
|                        | 1416          | 659        |
|                        | 1434          | 660        |
|                        | 1584          | 683        |
|                        | 1489          | 670        |
|                        | 1587          | 727        |
|                        | 1462          | 756        |
|                        | 1412          | 662        |
|                        | 1469          | 665        |
|                        | 1485          | 711        |
| Average                | 1477          | 686        |

### /reportcard

| API Response Time (ms) | Without Login | With Login |
| ---------------------- | ------------- | ---------- |
|                        | 1646          | 788        |
|                        | 1691          | 794        |
|                        | 1603          | 829        |
|                        | 1515          | 828        |
|                        | 1536          | 810        |
|                        | 2350          | 788        |
|                        | 1545          | 814        |
|                        | 1515          | 795        |
|                        | 1507          | 823        |
|                        | 1669          | 865        |
| Average                | 1658          | 813        |

### /schedule

| API Response Time (ms) | Without Login | With Login |
| ---------------------- | ------------- | ---------- |
|                        | 863           | 144        |
|                        | 911           | 109        |
|                        | 941           | 124        |
|                        | 966           | 129        |
|                        | 798           | 118        |
|                        | 816           | 123        |
|                        | 847           | 187        |
|                        | 921           | 128        |
|                        | 872           | 127        |
|                        | 969           | 180        |
| Average                | 890           | 137        |

### /transcript

| API Response Time (ms) | Without Login | With Login |
| ---------------------- | ------------- | ---------- |
|                        | 788           | 103        |
|                        | 823           | 98         |
|                        | 772           | 151        |
|                        | 914           | 107        |
|                        | 833           | 96         |
|                        | 837           | 96         |
|                        | 884           | 96         |
|                        | 891           | 98         |
|                        | 790           | 106        |
|                        | 812           | 113        |
| Average                | 834           | 106        |

## Tips

- Always POST to the `/login` endpoint before any subsequent requests, as it significantly boosts response times (see [Performance](#performance))
- Read the [documentation](#api-docs) to see if any parameters are avaliable in the body which might suit the use case

## Credits

- The packages mentioned in the [How It Works](#how-it-works) section
- The [template](https://github.com/create-go-app/fiber-go-template) the package was based on
- The [Frisco ISD API](https://github.com/SumitNalavade/FriscoISDHACAPI) by Sumit Nalavade which gave inspiration on choosing Go and how to gather the information
- The [HACify](https://github.com/Threqt1/HACify) repo for providing the base logic for parsing the raw HTML
