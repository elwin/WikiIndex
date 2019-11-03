# Wiki Index
This project has been created to play around with the Wikipedia dataset.
The goal is to calculate the distance between two pages, that is, how many clicks are required
to get from one page to the other by only using the links on each page.

First you should download the dataset, e.g. from [here](https://meta.wikimedia.org/wiki/Data_dump_torrents). Next, you can build and run the index:

```bash
go build
./WikiIndex -i ~/Downloads/simplewiki-20170820-pages-meta-current.xml.bz2 
```

After that, you can reach the web interface with [http://localhost:8080](http://localhost:8080).