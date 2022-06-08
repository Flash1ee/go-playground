# Докер и аналоги

## История возникновения докера
До появления контейнеризации существовали несколько способов изоляции приложений
### Chroot

Chroot - это системный вызов Unix-подобных ОС, позволяющий изменить корневой каталог в системе.

Этот механизм позволяет ограничить доступ программы к файлам в рамках одного каталога. 

#### Минусы
- если программа получит права супер-пользователя, она может сменить chroot и получить доступ к системе;
- нет ограничения на потребление ресурсов (память,  cpu);
- нет ограничений на доступ в сеть.

### Запуск экземпляров ОС на уровне исходной ОС

Метод виртуализации, при котором ядро операционной системы поддерживает несколько изолированных экземпляров пространства пользователя вместо одного. 
[[1]](#source-1)

#### Плюсы
- Обеспечивается изоляция по дисковому пространству, ресурсам, сети.
- Обеспечивается изоляция по ресурсам.

#### Минусы
- Необходимость поддерживать полноценную ОС внутри контейнера (большие расходы по ресурсам).
- Уязвимость ядра ОС (вероятность выйти из контейнера через уязвимость).


Плюсы и минусы этих способов сподвигли на создание контейнерной виртуализации.
Вместо запуска полноценной операционной системы в контейнере (с системой инициализации, пакетным менеджером и т.п.) можно запускать сразу же приложения, главное — обеспечить приложениям такую возможность (наличие необходимых библиотек и прочих файлов). 

Ниже рассмотрены релизации контейнерной виртуализации.

## Docker

Docker —  это открытая платформа для разработки, доставки и запуска приложений. Состоит из утилиты командной строки docker, которая вызывает одноименный сервис (сервис является потенциальной единой точкой отказа) и требует права доступа root. По умолчанию использует в качестве Container Runtime runc. Все файлы Docker (образы, контейнеры и др.) по умолчанию хранятся в каталоге /var/lib/docker.  [[3]](#source-3).

![image](https://user-images.githubusercontent.com/55987935/169495828-a5dcee3c-807b-4328-a71a-8b7bb264b63c.png)

### Архитектура

![image](https://user-images.githubusercontent.com/55987935/169496310-c3436c05-9270-4100-be61-18b137a0789c.png)

### Хранение данных

При запуске контейнер получает доступ на чтение ко всем слоям образа, а также создает свой исполняемый слой с возможностью создавать, обновлять и удалять файлы. Все эти изменения не будут видны для файловой системы хоста и других контейнеров, даже если они используют тот же базовый образ. При удалении контейнера все измененные данные также будут удалены.  
Чтобы расшарить между контйенерами данные существует два способа:

[named volumes](https://docs.docker.com/get-started/05_persisting_data/) – именованные тома хранения данных
Позволяет сохранять данные в именованный том, который располагается в каталоге в /var/lib/docker/volumes и не удаляется при удалении контейнера. Том может быть подключен к нескольким контейнерам

[bind mount](https://docs.docker.com/get-started/06_bind_mounts/) – монтирование каталога с хоста
Позволяет монтировать файл или каталог с хоста в контейнер. На практике используется для проброса конфигурационных файлов или каталога БД внутрь контейнера
```bash 
# named volume
docker run --detach --name jenkins --publish 80:8080 --volume=jenkins_home:/var/jenkins_home/ jenkins/jenkins:lts-jdk11 # запуск контейнера jenkins с подключением каталога /var/jenkins_home как тома jenkins_home
docker volume ls # просмотр томов
docker volume prune # удаление неиспользуемых томов и очистка диска. Для удаления тома все контейнеры, в которых он подключен, должны быть остановлены и удалены
 
# bind mount
# запуск контейнера node-exporter с монтированием каталогов внутрь контейнера в режиме read only: /proc хоста прокидывается в /host/proc:ro внутрь контейнера, /sys - в /host/sys:ro, а / - в /rootfs:ro
docker run \
-p 9100:9100 \
-v "/proc:/host/proc:ro" \
-v "/sys:/host/sys:ro" \
-v "/:/rootfs:ro" \
--name node-exporter prom/node-exporter:v1.1.2
```

<p><strong>Когда использовать тома, а когда монтирование с хоста</strong></p><br/>
<div class="scrollable-table"><table>
<thead>
<tr>
<th><u>Volume</u></th>
<th><u>Bind mount</u></th>
</tr>
</thead>
<tbody>
<tr>
<td>Просто расшарить данные между контейнерами.</td>
<td>Пробросить конфигурацию с хоста в контейнер.</td>
</tr>
<tr>
<td>У хоста нет нужной структуры каталогов.</td>
<td>Расшарить исходники и/или уже собранные приложения.</td>
</tr>
<tr>
<td>Данные лучше хранить не локально (а в облаке, например).</td>
<td>Есть стабильная структура каталогов и файлов, которую нужно расшарить между контейнерами.</td>
</tr>
</tbody>
</table></div><br/>

Подробнее: [[4]](#source-4)
## Podman

Podman – это инструмент с открытым исходным кодом для поиска, сборки, передачи и запуска приложений. Является утилитой командной строки с аналогичными docker командами, однако не требует дополнительный сервис для работы и может работать без прав доступа root. По умолчанию использует в качестве Container Runtime crun (ранее runc).

Возможность работать с контейнерами без прав `root` приводит к нескольким особенностям:

- все файлы Podman (образы, контейнеры и др.) пользователей с правами доступа `root` хранятся в каталоге `/var/lib/containers`, без прав доступа `root` – в `~/.local/share/containers`

- пользователи без root прав по умолчанию не могут использовать привилегированные порты и полноценно использовать некоторые команды

**Основные команды** Podman аналогичны docker, есть и приятные доработки (например, ключ --all для команд start, stop, rm, rmi).  
```bash
# справочная информация
podman --help # список доступных команд
podman <command> --help # информация по команде
 
# работа с образами
podman search nginx # поиск образов по ключевому слову nginx
 
podman pull ubuntu # скачивание последней версии (тег по умолчанию latest) официального образа ubuntu (издатель не указывается) из репозитория по умолчанию docker.io/library
podman pull quay.io/bitnami/nginx:latest # скачивание последней версии образа nginx от издателя bitnami из репозитория quay.io/bitnami
podman pull docker.io/library/ubuntu:18.04 # скачивание из репозитория docker.io официального образа ubuntu с тегом 18.04
 
podman images # просмотр локальных образов
 
podman rmi <image_name>:<tag> # удаление образа. Вместо <image_name>:<tag> можно указать <image_id>. Для удаления образа все контейнеры на его основе должны быть как минимум остановлены
podman rmi --all # удаление всех образов
 
# работа с контейнерами
podman run hello-world # Hello, world! в мире контейнеров
podman run -it ubuntu bash # запуск контейнера ubuntu и выполнение команды bash в интерактивном режиме
podman run --detach --name nginx --publish 9090:8080 quay.io/bitnami/nginx:1.20.2 # запуск контейнера nginx с отображением (маппингом) порта 9090 хоста на порт 8080 внутрь контейнера
podman run --detach --name mongodb docker.io/library/mongo:4.4.10 # запуск контейнера mongodb с именем mongodb в фоновом режиме. Данные будут удалены при удалении контейнера!
 
podman ps # просмотр запущенных контейнеров
podman ps -a # просмотр всех контейнеров (в том числе остановленных)
podman stats --no-stream # просмотр статистики. Если у пользователя нет прав доступа root, то необходимо переключиться на cgroups v2
 
podman create alpine # создание контейнера из образа alpine
 
podman start <container_name> # запуск созданного контейнера. Вместо <container_name> можно указать <container_id>
podman start --all # запуск всех созданных контейнеров
 
podman stop <container_name> # остановка контейнера. Вместо <container_name> можно указать <container_id>
podman stop --all # остановка всех контейнеров
 
podman rm <container_name> # удаление контейнера. Вместо <container_name> можно указать <container_id>
podman rm --all # удаление всех контейнеров
 
# система
podman system info # общая информация о системе
podman system df # занятое место на диске
podman system prune -af # удаление неиспользуемых данных и очистка диска
```

**Хранение данных** аналогично docker - присутствуют `named volumes` и `bind mount`.  
```bash
# справочная информация
podman <command> --help
 
# named volume
podman run --detach --name jenkins --publish 8080:8080 --volume=jenkins_home:/var/jenkins_home/ docker.io/jenkins/jenkins:lts-jdk11 # запуск контейнера jenkins с подключением каталога /var/jenkins_home как тома jenkins_home
podman volume ls # просмотр томов
podman volume prune # удаление неиспользуемых томов и очистка диска. Для удаления тома все контейнеры, в которых он подключен, должны быть остановлены и удалены
 
# bind mount
# запуск контейнера node-exporter с монтированием каталогов внутрь контейнера в режиме ro (read only)
podman run \
-p 9100:9100 \
-v "/proc:/host/proc:ro" \
-v "/sys:/host/sys:ro" \
-v "/:/rootfs:ro" \
--name node-exporter docker.io/prom/node-exporter:v1.1.2
```

**Containerfile** - аналог **DockerFile**.  

**Podman Compose** - аналог **Docker Compose**.

То, чего нет в docker - **Podman Pod**.

Podman Pod – это группа из одного или нескольких контейнеров с общим хранилищем и сетевыми ресурсами, а также спецификацией для запуска контейнеров. Концепция подов появилась и реализуется в Kubernetes [[5]](#source-5)

## Kata

Kata Containers — безопасная среда выполнения (runtime) контейнеров на основе облегченных виртуальных машин. Работа с ними происходит так же как и с другими контейнерами, но дополнительно имеется более надежная изоляция с использованием технологии виртуализации оборудования.  

Основные возможности:

- Работа с отдельным ядром, таким образом обеспечивается изоляция сети, памяти и операций ввода-вывода, есть возможность принудительного использования аппаратной изоляции на основе расширений виртуализации.
- Поддержка промышленных стандартов, включая OCI (формат контейнеров), Kubernetes CRI.
- Стабильная производительность обычных контейнеров Linux, повышение изоляции без накладных расходов, влияющих на производительность обычных виртуальных машин.
- Устранение необходимости запуска контейнеров внутри полноценных виртуальных машин, типовые интерфейсы упрощают интеграцию и запуск.

Источники:
1. <a id="source-1"></a> [Контейнеризация на уровне ОС](https://inlnk.ru/voyQA0)
2. <a id="source-2"></a> [Docker и все, все, все](https://habr.com/ru/company/southbridge/blog/512246/)
3. <a id="source-3"></a> [Основы контейнеризации](https://habr.com/ru/post/659049/)
4. <a id="source-4"></a> [Хранение данных Docker](https://habr.com/ru/company/southbridge/blog/534334/)
5. <a id="source-5"></a> [Podman Pod](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/building_running_and_managing_containers/assembly_working-with-pods_building-running-and-managing-containers)