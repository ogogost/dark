from PIL import Image
import cv2

im = Image.open('ss9.png')
im_crop = im.crop((250, 950, 1500, 1050))
im_crop.save('ss9task.png', quality=95)

image1 = cv2.imread('ss9task.png')
img = cv2.cvtColor(image1, cv2.COLOR_BGR2GRAY)
cv2.imshow('img', img)
ret, thresh1 = cv2.threshold(img, 160, 255, cv2.THRESH_BINARY)
cv2.imshow('thresh', thresh1)
cv2.imwrite('ss9task2.png',thresh1)
cv2.waitKey(0)