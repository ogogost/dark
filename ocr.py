import easyocr
# import numpy
import time

start = time.time() ## точка отсчета времени

reader = easyocr.Reader(['ru','en'], gpu=True)
result = reader.readtext('ss10_task.png', detail=0)

# file = open('ocr.txt', 'w')
# file.write(str(result))
# file.close()
# print(type(result))
print(result)
# print(sorted(result[0]))

end = time.time() - start ## собственно время работы программы

print(end) ## вывод времени