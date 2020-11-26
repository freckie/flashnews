import requests
from bs4 import BeautifulSoup

url = 'http://www.beyondpost.co.kr/list.php?ct=g0000'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='l2d').find('ul')
items = wrapper.find_all('li', recursive=False)
# print(items)

# exit(1)
req2 = requests.get('http://cnews.beyondpost.co.kr/view.php?ud=20201126171431341446a9e4dd7f_30')
bs2 = BeautifulSoup(req2.text, 'lxml')
print(bs2)

exit(1)

for item in items[0:3]:
    a_tag = item.find('a', class_='dyn std')
    href = 'https://dealsite.co.kr' + a_tag['href'].strip()
    title = a_tag.get_text().strip()
    date = item.find("div", class_='pubdate').get_text().strip()
    contents = item.find('span', class_='sneakpeek').get_text().strip()

    print('=================')
    print(title)
    print(date)
    print(href)
    print(contents)
