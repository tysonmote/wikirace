# wikirace

Given two Wikipedia pages, wikirace will try to find the shortest path from the
first page to the second page by only clicking on links in the main page
content.

## Limitations

wikirace adheres to [Wikipedia's published best practices for
bots][best_practices]. The main limitation here is:

> Do not make multi-threaded requests. Wait for one server request to complete
> before beginning another.

[best_practices]: https://en.wikipedia.org/wiki/Wikipedia:Creating_a_bot#Bot_best_practices
