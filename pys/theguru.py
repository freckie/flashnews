import requests
from bs4 import BeautifulSoup

url = 'https://www.theguru.co.kr/news/article_list_all.html'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='art_list_all')
items = wrapper.find_all('li', recursive=False)
# print(items)

# exit(1)
req2 = requests.get('https://www.theguru.co.kr/news/article.html?no=16579')
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
