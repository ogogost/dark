# import the necessary packages
import pytesseract
from PIL import Image

image = Image.open("ss1.png")
string = pytesseract.image_to_string(image, lang='rus')

print(string)