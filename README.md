Your client asks you to produce an "org chart" of their statically generated Website's navigation paths.

You can use the [Google sheets chart feature](https://support.google.com/docs/answer/9146871?hl=en) to draw, as long as the chart looks similar to what the client expects here:

<img src="https://s.natalian.org/2021-04-26/sheets.png">

Notice that a child can only have one parent, so you need to calculate the shortest path to the root element (aka "Home").

Your coding challenge is to create the same chart programatically for the
[static html](https://s.natalian.org/2021-04-26/site.zip) as well as their Website at https://simple.goserverless.sg/

Example output:

    generate-site-structure [index.html | https://simple.goserverless.sg/]
    Bread,Products
    Jam,Products
    Products,Home
    Privacy Policy,About Us
    Sustainability statement,About Us
    About Us,Home