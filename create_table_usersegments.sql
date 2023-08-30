create table UserSegments
(
user_id integer,
slug character varying (255) references segments (slug) ON DELETE CASCADE, --каскадное удаление, чтобы строки в зависимой таблице удалялись сами при удалении ключей в основной
primary key (user_id, slug) --составной ключ, чтобы каждая комбинация пользователь+сегмент была уникальной
)