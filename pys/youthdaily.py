import requests
from bs4 import BeautifulSoup

url = 'https://www.youthdaily.co.kr/news/article_list_all.html'
req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('ul', class_='art_list_all')
items = wrapper.find_all('li', recursive=False)

for item in items[:15]:
    a_tag = item.find('a')
    
    title = a_tag.find('h2').get_text().strip()
    url = 'https://www.youthdaily.co.kr' + a_tag['href']
    date = item.find('li', class_='date').get_text().strip()

    req2 = requests.get(url)
    bs2 = BeautifulSoup(req2.text, 'lxml')
    wrapper2 = bs2.find('div', id='news_body_area')

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])

    print(title)
    print(url)
    print(date)
    print(contents)