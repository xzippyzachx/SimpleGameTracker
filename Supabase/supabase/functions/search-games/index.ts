import "jsr:@supabase/functions-js/edge-runtime.d.ts"

import { corsHeaders } from '../_shared/cors.ts'
import { matchSorter } from 'npm:match-sorter'

Deno.serve(async (req) => {

  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders })
  }

  const { name } = await req.json()

  const IGDBAuthResponse = await fetch(`https://id.twitch.tv/oauth2/token?client_id=${Deno.env.get('CLIENT_ID')}&client_secret=${Deno.env.get('CLIENT_SECRET')}&grant_type=client_credentials`, {
    method: 'POST',
  });
  const IGDBAuthData = await IGDBAuthResponse.json();

  const body = `search "${name}"; limit 100; fields name, cover.image_id; where category = (0,4,8,9,11) & version_parent = null & cover != null;`

  const headers = new Headers({
    "Client-ID": Deno.env.get('CLIENT_ID'),
    "Authorization": `Bearer ${IGDBAuthData.access_token}`
  })

  const IGDBSearchResponse = await fetch(`https://api.igdb.com/v4/games`, {
    method: 'POST',
    body: body,
    headers: headers,
  });
  const IGDBSearchData = await IGDBSearchResponse.json();

  let sortedIGDBSearchData = matchSorter(IGDBSearchData, name, {keys: ['name']})

  sortedIGDBSearchData.forEach(game => {
    let imageId = game.cover.image_id
    game.cover = `https://images.igdb.com/igdb/image/upload/t_cover_big/${imageId}.jpg`
  });

  return new Response(
    JSON.stringify(sortedIGDBSearchData),
    { status: 200, headers: { ...corsHeaders, "Content-Type": "application/json" } },
  )
})
