import requests
from bs4 import BeautifulSoup

url = 'http://www.doctorsnews.co.kr/news/articleList.html?view_type=sm'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('section', class_='article-list-content')
items = wrapper.find_all('div', class_='list-block')

for item in items[:1]:
    div = item.find('div', class_='list-titles')
    a_tag = div.find('a')
    title = a_tag.get_text().strip()
    href = 'http://www.doctorsnews.co.kr' + a_tag['href']
    date = item.find('div', class_='list-dated').get_text().split(' | ')[2].strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', id='article-view-content-div')

    contents = ''.join([it.get_text().strip() for it in wrapper2.find_all('p')])
    print(title, href, date, contents)