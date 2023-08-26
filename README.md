# ltvpredictor

The repo contains the 2 versions of the requested utility.
- ltvpredictor_base - the first version, containing extended functionality such as:
  - full many to many relationship management, with Country and Campaigns rntities holding information about Analytics (links to objects, so we don't double the data)
  - with mentioned above we can extend functionality where needed, making cross references, seeing how certain campaigns behave in different countries, and such
  - Ability to budl graphs on predicted data. I predict all the days up to day 60, and use go-echart to visualize the data. charts are stored in the separate folder

- ltvpredictor_light - the stripped down version of the ltvpredictor_base
  - no many to many relations, we just gather the data that we need, depending on the values of the aggregate glag passed to the untility
  - country and campaign are now single entity, having name, analytics and usercount
  - this allows to save memory while parsing the data, however makes the code less readable and understandable

Both applications function in a similar fashion: check the input for validity, init the parser basing on detected file type, parse the files into corresponding structures, perform prediction on the parsed data.
Each prediction is done in a separate routine, results are passed back to main app where they are printed to the console.

All the information about how the utility functions can be obtained with: **go run . -h**

Running with **go run .** starts the app with a default values.
