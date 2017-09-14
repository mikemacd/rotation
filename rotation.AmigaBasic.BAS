'Amiga Basic
CLEAR ,130000&
DEFINT a,b,c,v,t
DEFDBL a1,b1,c1,a2,b2,c2,a3,b3,c3,th
PI=3.14
INPUT "Data File:";q$
OPEN q$ FOR INPUT AS 1
INPUT #1,v      'V is the number of verticies
DIM a(v,3)      'Dimensional array A holds the verticies of the object
DIM t(v,3)      'Temporary Dimensional array assoc. with A
DIM u(v,3)
FOR a=0 TO v-1
  INPUT #1,a(a,1),a(a,2),a(a,3)
NEXT
INPUT #1,t      'T is the number of triangle
DIM b(t,3)      'Dimensional array B describes the triangles
DIM theta(t)    'Array theta is for the angle between the light and the normal
FOR a=0 TO t-1
  INPUT #1,b(a,1),b(a,2),b(a,3)
NEXT
CLOSE 1

INPUT "XZ rotation factor:";zx
INPUT "XY rotation factor:";yx
INPUT "YZ rotation factor:";zy
INPUT "Magnification:";mag
CLS
l1=0
l2=0
l3=1
SCREEN 1,640,200,4,2
WINDOW 2,"Rotations",,0,1
WINDOW OUTPUT 2
PALETTE  0,0,0,0
PALETTE  1,.125 ,.125 ,.125
PALETTE  2,.1875,.1875,.1875
PALETTE  3,.25  ,.25  ,.25
PALETTE  4,.3125,.3125,.3125
PALETTE  5,.375 ,.375 ,.375
PALETTE  6,.4375,.4375,.4375
PALETTE  7,.5   ,.5   ,.5
PALETTE  8,.5625,.5625,.5625
PALETTE  9,.625 ,.625 ,.625
PALETTE 10,.6875,.6875,.6875
PALETTE 11,.75  ,.75  ,.75
PALETTE 12,.8125,.8125,.8125
PALETTE 13,.875 ,.875 ,.875
PALETTE 14,.9375,.9375,.9375
PALETTE 15,1,1,1

top:
  xz=xz+zx
  xy=xy+yx
  yz=yz+zy

  FOR a=0 TO v-1
    t(a,1)=a(a,1)*COS(xy)-a(a,2)*SIN(xy)
    t(a,2)=a(a,1)*SIN(xy)+a(a,2)*COS(xy)
    t(a,3)=a(a,3)

    u(a,1)=t(a,1)*COS(xz)-t(a,3)*SIN(xz)
    u(a,2)=t(a,2)
    u(a,3)=t(a,1)*SIN(xz)+t(a,3)*COS(xz)

    t(a,1)=u(a,1)
    t(a,2)=u(a,2)*COS(yz)-u(a,3)*SIN(yz)
    t(a,3)=u(a,2)*SIN(yz)+u(a,3)*COS(yz)

    u(a,1)=320 +t(a,1)*mag*2
    u(a,2)=100 +t(a,2)*mag
    u(a,3)=     t(a,3)*mag
  NEXT
  FOR a=0 TO t-1
    a1=u(b(a,2),1)-u(b(a,1),1)
    b1=u(b(a,2),2)-u(b(a,1),2)
    c1=u(b(a,2),3)-u(b(a,1),3)
    a2=u(b(a,3),1)-u(b(a,2),1)
    b2=u(b(a,3),2)-u(b(a,2),2)
    c2=u(b(a,3),3)-u(b(a,2),3)
    a3=b1*c2-c1*b2
    b3=c1*a2-a1*c2
    c3=a1*b2-b1*a2
    theta(a)= ((a3*l1)+(b3*l2)+(c3*l3)) /
              (((a3^2+b3^2+c3^2)^.5)*((l1^2+l2^2+l3^2)^.5))
  NEXT
  CLS
  FOR a=0 TO t-1
    q$=INKEY$
    IF q$<>"" THEN GOSUB keys
    IF theta(a)<0 THEN nxt
    col=CINT(theta(a)*15)
    COLOR col
    AREA (u(b(a,1),1),u(b(a,1),2))
    AREA (u(b(a,2),1),u(b(a,2),2))
    AREA (u(b(a,3),1),u(b(a,3),2))
    AREAFILL
nxt:
  NEXT
  GOTO top

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
