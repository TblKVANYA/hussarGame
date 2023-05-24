# CardGame Hussar
Несколько лет назад мой отец научил меня карточной игре, придуманной им в детстве. Игра мне понравилась, и вот теперь я решил поделиться ей с миром. В текущем виде она рассчитана на двоих-четверых игроков. В будущем возможно расширение на бОльшее количество игроков.

## Правила
Пусть в игре участвуют N игроков, играющих стандартной колодой на 36 карт.

**Как победить?**

Набрать больше очков, чем соперники.

**За что начисляются очки?**

Игра состоит из раундов. В каждом из них игрок зарабатывает или теряет очки. Перед началом каждого раунда определяется козырь, игрок получает карты, оценивает их и делает ставку, сколько *раздач* он выиграет. Пусть ставка игрока - Х. В конце раунда сраниваются предсказание игрока и количество *раздач* выигранных в реальности - Y. Если X = Y, игрок получает 10\*X очков. Если X < Y, игрок получает Y очков. Если X > Y, игрок теряет 10\*X очков.

**Что такое раздача и как ее выиграть?** 

Один из игроков за столом скидывает любую свою карту. Вслед за ним это делают и остальные по часовой стрелке. Чтобы выиграть раздачу, нужно перебить ведущую карту. Сделать это можно либо козырем, либо картой той же масти, но выше ранком. Тогда сделавшая это карта становится ведущей. "Слить" карту нельзя (т.е. нельзя кинуть козырного туза так, чтобы он не стал ведущей картой). Обладатель ведещей карты становится победителем раздачи.

**Все раунды одинаковые? Кто ходит первым?**

Сначала идут несколько раундов, отличающихся только количеством карт. В начале должно быть N раундов, в которох игроки получают по одной карте. Затем по одному раунду, в котором игроки получают от *2* до *(36/N)-1* карт. Затем N раундов, в которых игроки получают *36/N* карт, а после количество карт постепенно уменьшается и эта часть заканчивается как начиналась - N раундов по одной карте на игрока.

Зачем такие странности? Тут добавляется новое правило - нельзя N раундов подряд ставить, что ты выиграешь 0 раздач.

В конце идут три особых раунда - темный, бескозырный и золотой. Во всех них игроки получают на руки по 12 карт. В темном ставки идут до того, как игроки увидят свои карты. В бескозырном нет козырей. Количество очков, полученных в золотом раунде, утраивается.

По поводу порядка ходов - в начале игры случайно выбирается дилер. В каждом раунде дилер первым ставку и начинает первую раздачу. Дальнейшие раздачи начинает игрок, выигравший предыдущую раздачу этого раунда. После окончания раунда роль дилера переходит следующему по часовой стрелке.

**Окей, примерно понятно. Есть ли еще нюансы?** 

Да, есть. Крестовый валет - магическая карта, способная принимать любые значения из игровой колоды и даже больше. Он может стать пятеркой или супериором(мне пришлось выдумать что-то значимее туза) любой масти, в зависимости от желаний того, кто его кинул. Значение нужно определить перед сбросом карты.

## Установка и запуск
Проект состоит из двух частей - сервера и консольного клиента, лежащих в соответсвующих папках.

1. [Установите язык Go](https://go.dev/doc/install), настройте переменную GOPATH
2. [Запустите сервер](/server/)
3. [Запустите клиенты.](/client/) Первый, кто подключился, выбирает количество игроков.
4. Наслаждайтесь:)

## Предостережение
Клиент и сервер рассчитаны на игру по локальной сети