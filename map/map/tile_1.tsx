<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.2" tiledversion="1.2.1" name="tile_1" tilewidth="64" tileheight="64" tilecount="12" columns="3">
 <image source="tile_1.png" width="192" height="256"/>
 <tile id="1">
  <properties>
   <property name="type" value="barrier"/>
  </properties>
  <objectgroup draworder="index">
   <object id="2" x="62.6667" y="0.66667">
    <properties>
     <property name="boundary_type" value="barrier"/>
    </properties>
    <polyline points="2,-3 -65.6667,-6.33333 -65.6667,63.6667 3.66667,68.3333 5.33333,1.33333"/>
   </object>
  </objectgroup>
 </tile>
 <tile id="4">
  <properties>
   <property name="type" value="barrier"/>
  </properties>
  <objectgroup draworder="index">
   <object id="2" x="63.3333" y="-0.666667">
    <properties>
     <property name="boundary_type" value="barrier"/>
    </properties>
    <polyline points="0,0 -64,0 -64,64.6667 0.666667,64.6667 1.33333,10.6667"/>
   </object>
  </objectgroup>
 </tile>
 <tile id="5">
  <properties>
   <property name="type" value="barrier"/>
  </properties>
  <objectgroup draworder="index">
   <object id="5" x="1" y="3">
    <properties>
     <property name="boundary_type" value="barrier"/>
    </properties>
    <polyline points="0,0 64,-2 63,55 -1,61 -2,9"/>
   </object>
  </objectgroup>
 </tile>
 <tile id="7">
  <properties>
   <property name="type" value="lowScoreTreasure"/>
  </properties>
  <objectgroup draworder="index">
   <object id="4" x="32" y="8">
    <polyline points="-31,-10 -33,18 -18,46 17,49 32,24 6,-10"/>
   </object>
  </objectgroup>
 </tile>
 <tile id="8">
  <properties>
   <property name="type" value="highScoreTreasure"/>
  </properties>
  <objectgroup draworder="index">
   <object id="1" x="32" y="10">
    <polyline points="-90,-119 -216,17 -97,181 177,168 234,3 55,-118"/>
   </object>
  </objectgroup>
 </tile>
 <tile id="9">
  <properties>
   <property name="type" value="trap"/>
  </properties>
  <objectgroup draworder="index">
   <object id="1" x="19" y="-52">
    <polyline points="-115,-135 -351,67 -187,314 196,312 340,65 116,-137"/>
   </object>
  </objectgroup>
 </tile>
 <tile id="10">
  <properties>
   <property name="type" value="speedShoes"/>
  </properties>
  <objectgroup draworder="index">
   <object id="1" x="-18" y="-18">
    <polyline points="0,0 -11,82 85,120 126,50 90,-20 33,-9 31,-8 33,-7"/>
   </object>
  </objectgroup>
 </tile>
</tileset>
