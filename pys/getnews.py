import requests
from bs4 import BeautifulSoup

url = 'http://www.getnews.co.kr/list.php?ct=g0000'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='l2d').find('ul')
items = wrapper.find_all('li')

for item in items[0:10]:
    href = item.find('a')['href'].strip()
    title = item.find('span', class_='w1 elip1').get_text().strip()
    date = item.find("span", class_='e2').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    
    wrapper2 = bs2.find('div', class_='vcon_in articleContent')
    contents = ''.join([it.strip() for it in wrapper2.find_all(recursive=False, text=True)])

    print('=================')
    print(title)
    print(date)
    print(href)
    print(contents)