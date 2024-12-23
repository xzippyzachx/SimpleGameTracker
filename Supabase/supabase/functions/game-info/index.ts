// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

// Setup type definitions for built-in Supabase Runtime APIs
import "jsr:@supabase/functions-js/edge-runtime.d.ts"

import { corsHeaders } from '../_shared/cors.ts'

Deno.serve(async (req) => {
  
  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders })
  }
  
  const { gameId } = await req.json()

  const IGDBAuthResponse = await fetch(`https://id.twitch.tv/oauth2/token?client_id=${Deno.env.get('CLIENT_ID')}&client_secret=${Deno.env.get('CLIENT_SECRET')}&grant_type=client_credentials`, {
    method: 'POST',
  });
  const IGDBAuthData = await IGDBAuthResponse.json();

  const body = `where id = ${gameId}; fields name, cover.image_id, summary;`

  const headers = new Headers({
    "Client-ID": Deno.env.get('CLIENT_ID'),
    "Authorization": `Bearer ${IGDBAuthData.access_token}`
  })

  const IGDBGameResponse = await fetch(`https://api.igdb.com/v4/games`, {
    method: 'POST',
    body: body,
    headers: headers,
  });
  let gameData = await IGDBGameResponse.json();

  if (gameData.length != 1)
  {
    return new Response(
      "Invalid gameID",
      { status: 400, headers: { ...corsHeaders, "Content-Type": "application/json" } },
    )
  }

  gameData = gameData[0]
  let imageId = gameData.cover.image_id
  gameData.cover = `https://images.igdb.com/igdb/image/upload/t_cover_big/${imageId}.jpg`

  return new Response(
    JSON.stringify(gameData),
    { status: 200, headers: { ...corsHeaders, "Content-Type": "application/json" } },
  )

})

