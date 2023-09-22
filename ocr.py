import easyocr
# import numpy

reader = easyocr.Reader(['ru','en'], gpu=False)
result = reader.readtext('ss9task2.png', detail=0)

file = open('ocr.txt', 'w')
file.write(str(result))
file.close()
