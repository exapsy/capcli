# Architectural Choices

## HTTP 1.1 vs HTTP 2.0 vs HTTP 3.0

HTTP 2.0 and HTTP 3.0 require TLS.

HTTP 3.0 is not yet widely supported either but I didn't even try in Capital.com's backend API.

The **problem** is that they require TLS which makes it useless for a simple CLI tool to
use HTTP 2 or HTTP 3. SSL/TLS, despite that it's a PITA (Pain in the ass) to maintain or to complain to the user
that they didn't use the correct configuration, the real reason is that
simply processing on every CLI command a whole chain of TLS configuration LOC (lines of code),
would probably not only make the performance gain insignificant from HTTP 2 or HTTP 3,
but also probably would be more performance heavy.

## Why not a simple CLI?

Guys and girls that want to use a GUI, are welcome. That's why.
Not everybody is a nerd.

## Why not have GUI instead of GUI in Terminal?

I'll work on it. I like the idea, but I only have so much time.
I love making tools for people that find them useful, and even more
when they're more accessible to everybody, and a GUI would certainly do that.

But first I'm working on making an infrastructure, then on the frontend.
