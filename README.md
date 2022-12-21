
Dashboard to display cpu and memory load in real-time using Golang, vanilla JS , chartjs library.

General features:

- Real-Time display CPU and Memory usage in %

Application design:

1. Golang as backend and sent Server Sent Event(SSE) to web client
2. Vanilla JS to read SSE
3. chartjs to display realtime chart

Installation (local):

1. docker build -t go-sse-demo-first .
2. docker run --rm -p 5000:5000 go-sse-demo-first

Installation (fly.io):

1. flyctl launch

Reference(s):

[Server-Sent Events with Go and React](https://articles.wesionary.team/server-sent-events-with-go-and-react-76df101a3efe)

[Infinite looping with vs without time.Sleep](https://stackoverflow.com/questions/55858835/infinite-looping-with-vs-without-time-sleep)

[Updating Charts](https://www.chartjs.org/docs/latest/developers/updates.html)

[Multi Axis Line Chart](https://www.chartjs.org/docs/latest/samples/line/multi-axis.html)
