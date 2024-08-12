# Kia Storm's Digi2al Submission

## Running the code

> I'm using `go version go1.22.5 darwin/arm64`.

1. `go run main.go` to run the code.

## What I did

> 1. Careful tool selection to satisfy the above requirements, you may opt to build everything from scratch or choose appropriate tools

I chose to Go + HTMX for simplicity and ease of use, given the requirements.

> 2. The login journey in pure HTML/CSS (you may use some javascript to mock the login system or use a http server with basic auth)

I chose a mock authentication system using vanilla Go.

> 3. Progressively enhanced version, which may or may not dynamically update the page

The `toggle-password` web component is progressively enhanced, only showing the toggle when JavaScript is available.

> 4. WCAG 2.1 conformance levels

I've tried to make the project as accessible as possible within reasonable time, but I'm sure there are areas for improvement and I would love to hear feedback on how else to improve it.

## What else would I do?

- Incorporate the learnings found here https://technology.blog.gov.uk/2021/04/19/simple-things-are-complicated-making-a-show-password-option/
  - More specifically, I would announce the change of state to screen readers when toggling between show/hide
- Documentation, tests, better organisation etc.

Hope you like it!

Kia
kia@kormsen.com
