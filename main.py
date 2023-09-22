import pyautogui
import win32gui
import cv2
from PIL import Image
def screenshot(window_title=None):
    if window_title:
        title_exists = win32gui.FindWindow(None, window_title)
        if title_exists:
            win32gui.SetForegroundWindow(title_exists)
            x, y, x1, y1 = win32gui.GetClientRect(title_exists)
            x, y = win32gui.ClientToScreen(title_exists, (x, y))
            x1, y1 = win32gui.ClientToScreen(title_exists, (x1 - x, y1 - y))
            im = pyautogui.screenshot('ss10.png', region=(x, y, x1, y1))
            return im
        else:
            print('Не существует окна с таким именем.')
    else:
        im = pyautogui.screenshot()
        return im

try:
    im = screenshot('Seekers Notes®: Тайны Дарквуда') #Пишешь название окна в точности до символа.
    im.show()
except:
    print('Не существует окна с таким именем.')

im = Image.open('ss10.png')
im_crop = im.crop((250, 950, 1500, 1050))
im_crop.save('ss10task.png', quality=95)

# img = cv2.imread('ss6.png', cv2.IMREAD_GRAYSCALE)
# cv2.imshow('game', img)
# cv2.waitKey(0)