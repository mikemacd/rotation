# rotation
This is an experiment I originally wrote around 1988 in Amiga Basic. It is a program which will load and draw a three dimensional object. The math needed to do this I wouldn't learn 'officially' until several years later in university. 

As a coding challenge I ported this program to C under linux/solaris using X11 around 1996.

Once again I have undertaken as an exercise to rewrite this application, this time (2017) in Go.

Since this uses `github.com/andlabs/ui` as the graphics engine there may need to be some prerequisites installed in order to build this such as MinGW on windows or gtk3 on linux. 

Apart from that running it is simply a matter of:

```
go build
./rotation -f icosahedron.dat
```

extended help information is available
```
./rotation -h
```
