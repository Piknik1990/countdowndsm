# Countdown DSM

<p align="center">
  <br>
  <img src="demo.gif" width="600" alt="CountdownDSM Demo">
  <br>
</p>

## Концепция

Таймер для выступления определённого перечня выступающих, строго регламентированных по времени.

Каждый участник проходит все заданные этапы выступления. Если выступающий успел раньше, можно нажать
`Space` или `Enter` и перейти на следующий этап. Если выступающий не успел - значит опоздал.

## Установка

Программа готова к работе на Linux, Mac и Windows. Cкачайте бинарный файл из [releases](https://github.com/Piknik1990/countdowndsm/releases).

Пример установки для Linux:

```shell
$ wget https://github.com/Piknik1990/countdowndsm/releases/download/v1.3.0/countdowndsm-linux-amd64
$ sudo cp countdowndsm-linux-amd64 /usr/bin/countdowndsm
$ sudo chmod 755 /usr/bin/countdowndsm
$ countdowndsm
 countdowndsm <pathtoconfig>

 Usage
        countdowndsm /path/to/config.yml
```

## Использование

Настройка работы приложения происходит через yaml-файл, путь до которого передаётся в виде аргумента

```sh
countdowndsm /path/to/config.yml
```

Файл конфигурации содержит следующие параметры:

* `persons` - перечень имён выступающих. Каждый из них будет проходить этапы выступление из `acts`
* `random` - флаг перемешивания персон. Если `false` - выступающие будут в случайном порядке; `true` - строго по списку
* `acts` - порядок этапов выступлений для каждой из персон выше:
  *  `name` - название этапа
  *  `time` - время этапа
* `counter` - добавить счётчик выступающих
* `next` - показывать следующего выступающего

## Горячие клавиши

* `Space` или `Enter`: Пропустить текущий этап
* `Esc` или `Ctrl+C`: Остановить работу программы
* `Tab`: Паузирование таймера

## Лицензия

[MIT](LICENSE)
