import requests
from bs4 import BeautifulSoup

url = 'http://www.wowtv.co.kr/NewsCenter/News/NewsList?subMenu=latest&menuSeq=459'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='contain-list-news')
items = wrapper.find_all('div', class_='article-news-list')

for item in items[:3]:
    div = item.find('div', class_='contian-news photo-right')
    a_tag = div.find('a')

    title = a_tag.find('p').find(text=True, recursive=False).strip()
    url = '' + a_tag['href']
    date = a_tag.find('span', class_='date').get_text().strip()

    req2 = requests.get(url)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='divNewsContent')

    contents = ''.join([it.strip() for it in wrapper2.find_all(text=True, recursive=False)])

    print(title)
    print(url)
    print(date)
    print(contents)