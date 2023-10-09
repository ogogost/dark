import easyocr
import time

start = time.time() ## точка отсчета времени

reader = easyocr.Reader(['ru'], gpu=True)

ready = time.time() ## точка отсчета времени

print('reader loaded: ' + str(round((ready - start), 2)) + ' seconds') ## вывод времени загрузки модели

result = reader.readtext('ss10_task.png', detail=0) # распознование текста с картинки
print(result)

result1 = reader.readtext('ss9.png', detail=0) # распознование текста с картинки
print(result1)

result2 = reader.readtext('ss8.png', detail=0) # распознование текста с картинки
print(result2)



end = time.time() - ready ## собственно время работы программы

print('ocr time: ' + str(round(end, 2)) + ' seconds') ## вывод времени