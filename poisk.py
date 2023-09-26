import cv2
import numpy as np

img = cv2.imread('ss2.png',1)
# fon = cv2.imread('fon2.png')
gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

items = {
         'skripka': ['imgs/bal_zal/skripka/1.png', 'imgs/bal_zal/skripka/2.png'],
         'horse': ['imgs/bal_zal/horse/1.png', 'imgs/bal_zal/horse/2.png'],
         'binokl': ['imgs/bal_zal/binokl/1.png', 'imgs/bal_zal/binokl/2.png'],
         'svitok': ['imgs/bal_zal/svitok/1.png', 'imgs/bal_zal/svitok/2.png'],
         'pes_chas': ['imgs/bal_zal/pes_chas/1.png', 'imgs/bal_zal/pes_chas/2.png'],
         'zerkalce': ['imgs/bal_zal/zerkalce/1.png', 'imgs/bal_zal/zerkalce/2.png'],
         }
templates = ['skripka', 'horse', 'binokl', 'svitok', 'zerkalce', 'pes_chas']

dlina_spiska_putei = 2

template_paths = []
for i in templates:
    for j in range(dlina_spiska_putei):
        template_paths.append(items[i][j])

print(len(template_paths))

for i in range(len(template_paths)):
   file_path = template_paths[i]
   template = cv2.imread(file_path,0)
   w,h = template.shape[0], template.shape[1]
   matched = cv2.matchTemplate(gray,template,cv2.TM_CCOEFF_NORMED)
   threshold = 0.8
   loc = np.where(matched >= threshold)
   for pt in zip(*loc[::-1]):
      cv2.rectangle(img, pt, (pt[0] + h, pt[1] + w), (0,255,255), 5)

cv2.imshow('Matched with Template',img)
cv2.imwrite('ss2_item.png', img)
cv2.waitKey(0)
