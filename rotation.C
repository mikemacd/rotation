/* Rotation - A program to rotate an object in three dimensions     */
/*                                                                  */
/* By: Michael A. MacDonald                 mikemacd@cyberspace.org */
/* Last update: September 19, 1996                                  */
/*                                                                  */
/*                                                                  */



/* include the X functions...*/
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
/* include the X library headers */
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xos.h>
#include <X11/Shell.h>
#include <X11/Xatom.h>

/* here are our X variables */
Display *dis;
int screen;
Window win;
GC gc;
XColor  rgb;

unsigned long black,white,colour;
int grey;

/* here are our X routines declared! */
void init_x();
void close_x();
void redraw();
void DrawTriangle();
void DrawTriangleM();
void DrawLineM();

#define MAX_VERTEX 500
#define MAX_TRIANGLES 500



extern unsigned sleep();

main (argc, argv)
    int argc;
    char *argv[];
{
       init_x();


        /* Define the light source location */
        float l1=0,l2=0,l3=10;

        /* Define the screen size */
        float   ScreenWidth=700, ScreenHeight=700;

        float   verticies[MAX_VERTEX][4];       /* verticies of triangles V(x,y,z) */
        float   verttemp1[MAX_VERTEX][4];
        float   verttemp2[MAX_VERTEX][4];
        int     triangles[MAX_TRIANGLES][4];    /* each triangle is a vector of verticies */
        float   theta[MAX_TRIANGLES];           /* angle between the normal of a triangle and the light source */

        float   xz=.1,zx=.1;                            /* rotation factors in the x-z plane */
        float   xy=.1,yx=.1;                            /* rotation factors in the x-y plane */
        float   yz=.1,zy=.1;                            /* rotation factors in the y-z plane */
        float   mag=1,delay=0;
        int     loops=0,a,v,t,wireframe=0;
        float   a1,b1,c1,a2,b2,c2,a3,b3,c3,th;

        FILE    *fpin;
        int     err;
        int     exit_cond=1;


        /* open and read the data file */
        if (argc < 2)
        {
                puts("Missing Data Filename.\n");
                exit(1);
        }

        if (argc>=3)
        {
                wireframe=atoi(argv[2]);
        }
        if (argc>=4)
        {
                delay = (100000 * atof(argv[3]));
        }

        if ((fpin=fopen(argv[1],"r")) == NULL)
        {
                fprintf(stderr,"Cannot read %s.",argv[1]);
                exit(1);
        }
        else
        {
//              if (!wireframe)
//                      PrepareColors();

                fscanf(fpin,"%d",&v);                   /* Read the number of verticies in the object */
                for (a=0; a<=v-1; a++)
                {
                        /* read a vertex */
                        fscanf(fpin,"%f %f %f",&verticies[a][1],&verticies[a][2],&verticies[a][3]);
                }

                fscanf(fpin,"%d",&t);                   /* Read the number of triangles in the object */
                for (a=0; a<=t-1; a++)
                {
                        /* read a triangle */
                        fscanf(fpin,"%d %d %d",&triangles[a][1],&triangles[a][2],&triangles[a][3]);
                }

                err=fclose(fpin);
        }

        /* read the rotation factors from stdin */
        printf("\nEnter XZ rotation factor: ");
                scanf("%f",&zx);
        printf("\nEnter XY rotation factor: ");
                scanf("%f",&yx);
        printf("\nEnter YZ rotation factor: ");
                scanf("%f",&zy);
        printf("\nEnter magnification factor: ");
                scanf("%f",&mag);
        while (exit_cond)
        {
                /* Change the rotation angle by all three rotation factors. */
                xz=xz+zx;
                xy=xy+yx;
                yz=yz+zy;

                if (xz>9.424778)
                        xz-=9.424778;
                if (xy>9.424778)
                        xy-=9.424778;
                if (yz>9.424778)
                        yz-=9.424778;

                for (a=0; a<=v-1;a++)
                {
                        /* Multiply each vector by the rotation matrix */
                        verttemp1[a][1]  =  (verticies[a][1]*cos(xy)  -  verticies[a][2]*sin(xy));
                        verttemp1[a][2]  =  (verticies[a][1]*sin(xy)  +  verticies[a][2]*cos(xy));
                        verttemp1[a][3]  =  (verticies[a][3]);

                        verttemp2[a][1]  =  (verttemp1[a][1]*cos(xz)  -  verttemp1[a][3]*sin(xz));
                        verttemp2[a][2]  =  (verttemp1[a][2]);
                        verttemp2[a][3]  =  (verttemp1[a][1]*sin(xz)  +  verttemp1[a][3]*cos(xz));

                        verttemp1[a][1]  =  (verttemp2[a][1]);
                        verttemp1[a][2]  =  (verttemp2[a][2]*cos(yz)  -  verttemp2[a][3]*sin(yz));
                        verttemp1[a][3]  =  (verttemp2[a][2]*sin(yz)  +  verttemp2[a][3]*cos(yz));

                        /* Scale the object by the magnification factor and adjust for the screen */
                        verttemp2[a][1]  =  ((ScreenWidth/2) +verttemp1[a][1]*mag);
                        verttemp2[a][2]  =  ((ScreenHeight/2)+verttemp1[a][2]*mag);
                        verttemp2[a][3]  =  (                 verttemp1[a][3]*mag);

                } /* End of rotate object */

                /* Calculate the angle between the normal of the triangle and the lightsource */
                /* This tells us how much light will "hit" the surface. */
                for (a=0; a<=t-1;a++)
                {

                        a1  =  verttemp2[triangles[a][2]][1]  -  verttemp2[triangles[a][1]][1];
                        b1  =  verttemp2[triangles[a][2]][2]  -  verttemp2[triangles[a][1]][2];
                        c1  =  verttemp2[triangles[a][2]][3]  -  verttemp2[triangles[a][1]][3];

                        a2  =  verttemp2[triangles[a][3]][1]  -  verttemp2[triangles[a][2]][1];
                        b2  =  verttemp2[triangles[a][3]][2]  -  verttemp2[triangles[a][2]][2];
                        c2  =  verttemp2[triangles[a][3]][3]  -  verttemp2[triangles[a][2]][3];

                        a3  =  b1*c2  -  c1*b2;
                        b3  =  c1*a2  -  a1*c2;
                        c3  =  a1*b2  -  b1*a2;


                        theta[a]  =   ((a3*l1) + (b3*l2) + (c3*l3))   /   (sqrt((a3*a3 + b3*b3 + c3*c3)) * sqrt((l1*l1 + l2*l2 + l3*l3)));
                } /* End of calculate normal */

                /*Delay screen clear to let image register*/
                for (a=0;a<=((int)(delay*10));a++)
                { /*Do nothing*/        }

                /*Clear the screen*/
                XClearWindow(dis, win);
                XFlush(dis);

/*
                loops++;
                SetPen((int)((float)loops/25));
*/

                for (a=0; a<=t-1; a++)
                {

/*                      ** GET KEY PRESS **
*/
                        if ((theta[a] >= 0) || (wireframe==1))
                        {
                                /* Draw the triangles only if light will reach their surface */
                                /* col=CINT(theta[a]*63)                                     */
                                /* We only need to use the x and y value of a point          */


//                              XFillPolygon(dis,win,gv,FillTriangle,3,1,0);

                                grey= (int)(theta[a]*65536);
                                rgb.red=grey;
                                rgb.green=grey;
                                rgb.blue=grey;
                                rgb.flags = DoRed | DoGreen | DoBlue;

if (!wireframe)                 DrawTriangle(   (int) (verttemp2[ triangles[a][1] ][1]),(int) (verttemp2[ triangles[a][1] ][2]),
                                                (int) (verttemp2[ triangles[a][2] ][1]),(int) (verttemp2[ triangles[a][2] ][2]),
                                                (int) (verttemp2[ triangles[a][3] ][1]),(int) (verttemp2[ triangles[a][3] ][2]),
                                                &rgb
                                );
if (wireframe)                  DrawTriangleM(  (int) (verttemp2[ triangles[a][1] ][1]),(int) (verttemp2[ triangles[a][1] ][2]),
                                                (int) (verttemp2[ triangles[a][2] ][1]),(int) (verttemp2[ triangles[a][2] ][2]),
                                                (int) (verttemp2[ triangles[a][3] ][1]),(int) (verttemp2[ triangles[a][3] ][2]),
                                                (( (theta[a]>=0) * ((int)(theta[a]*63)+47)) + (15*(theta[a]<0)) ));

                        } /*End of draw triangle */
                } /* End of draw object */
                XFlush(dis);
        } /* End of loop while exit_cond */
} /* End of main() */




void DrawTriangle(x1, y1, x2, y2, x3, y3, col)
int x1, y1, x2, y2, x3, y3,col;
{
         XPoint pts[6];
        XSetForeground(dis, gc,col);

        pts[0].x=x1;    pts[0].y=y1;

        pts[1].x=x2;    pts[1].y=y2;

        pts[2].x=x3;    pts[2].y=y3;

        XFillPolygon(dis,win,gc,pts,3,Convex,0);
}

void DrawTriangleM(x1, y1, x2, y2, x3, y3, col)
int x1, y1, x2, y2, x3, y3,col;
{
        XSetForeground(dis, gc,col);
        DrawLineM(x1, y1, x2, y2);
        DrawLineM(x2, y2, x3, y3);
        DrawLineM(x3, y3, x1, y1);
}

void DrawLineX(x1, y1, x2, y2)
int x1, y1, x2, y2;
{
        XDrawLine(dis, win, gc, x1, y1, x2, y2);
        XFlush(dis);
}

void DrawLineM(x1, y1, x2, y2)
int x1, y1, x2, y2;
{
        XDrawLine(dis, win, gc, x1, y1, x2, y2);
}


void init_x() {
/* get the colors black and white (see section for details) */
        dis=XOpenDisplay((char *)0);
        screen=DefaultScreen(dis);
        black=BlackPixel(dis,screen),
        white=WhitePixel(dis, screen);

 printf ("black:%ld\n", black);
 printf ("white:%ld\n", white);
//              OpenX("Rotations - Michael A. MacDonald",(int) ScreenWidth, (int) ScreenHeight);
        win=XCreateSimpleWindow(dis,DefaultRootWindow(dis),0,0,600, 600, 5,black, white);
        XSetStandardProperties(dis,win,"Rotations - Michael A. MacDonald","Hi",None,NULL,0,NULL);
         XSelectInput(dis, win, ExposureMask|ButtonPressMask|KeyPressMask);
        gc=XCreateGC(dis, win, 0,0);
        XSetBackground(dis,gc,white);
        XSetForeground(dis,gc,black);
        XClearWindow(dis, win);
        XMapRaised(dis, win);
};

void close_x() {
        XFreeGC(dis, gc);
        XDestroyWindow(dis,win);
        XCloseDisplay(dis);
        exit(1);
};

void redraw() {
        XClearWindow(dis, win);
};









/* Deal with key presses
keys:
  IF q$="7" THEN
    yx=yx+.025
  ELSEIF q$="8" THEN
    zx=zx+.025
  ELSEIF q$="9" THEN
    zy=zy+.025
  ELSEIF q$="1" THEN
    yx=yx-.025
  ELSEIF q$="2" THEN
    zx=zx-.025
  ELSEIF q$="3" THEN
    zy=zy-.025
  ELSEIF q$="4" THEN
    yx=0
  ELSEIF q$="5" THEN
    zx=0
  ELSEIF q$="6" THEN
    zy=0
  ELSEIF q$="(" THEN
    mag=mag-.05
  ELSEIF q$=")" THEN
    mag=mag+.05
  END IF
  IF mag=<0 THEN mag = .05
  q$=""
  RETURN
*/
