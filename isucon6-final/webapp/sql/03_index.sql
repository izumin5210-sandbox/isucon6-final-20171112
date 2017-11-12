create index
  index_strokes_on_room_id_and_id
on strokes (
  room_id,
  id
)

create index
  index_points_on_stroke_id_and_id
on strokes (
  stroke_id,
  id
)

create index
  index_tokens_on_csrf_token_and_created_at
on strokes (
  csrf_token,
  created_at
)
