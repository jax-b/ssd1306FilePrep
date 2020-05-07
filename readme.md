# SSD1306 File Prep
this application takes a image and converts it into raw Black and White Bytes for the SSD 1306 Display. This program will auto resize the image to fit it within the screen and convert it to black and white using a threshold value.



### Usage

| Command       | Type   | Description                                                                             |
| ------------- |:------:| ---------------------------------------------------------------------------------------:|
| -inputFile    | string | input image file. (Required)                                                            |
| -outputFile   | string | output .b file. (Required)                                                              |
| -displaySizeX | int    | sets the with of the screen (defaults to 128) (default 128)                             |
| -displaySizeY | int    | sets the Height of the screen (defaults to 64) (default 64)                             |
| -threshold    | int    | sets the threshold value for the conversion (defaults to 128, from 0-255) (default 128) |

