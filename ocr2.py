file = open('ocr.txt', 'r')
while True:
    # считываем строку
    line = file.readline()
    # прерываем цикл, если строка пустая
    if not line:
        break
    # выводим строку
    # print(line.strip())
    print(line)

# закрываем файл
file.close
