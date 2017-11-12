alter table rooms add column stroke_count int default 0;
update rooms set rooms.stroke_count = (select count(*) from strokes where strokes.room_id = rooms.id);
