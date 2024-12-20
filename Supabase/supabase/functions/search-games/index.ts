// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

// Setup type definitions for built-in Supabase Runtime APIs
import "jsr:@supabase/functions-js/edge-runtime.d.ts"

import { matchSorter } from 'npm:match-sorter'

Deno.serve(async (req) => {
  const { name } = await req.json()

  const IGDBAuthResponse = await fetch(`https://id.twitch.tv/oauth2/token?client_id=${Deno.env.get('CLIENT_ID')}&client_secret=${Deno.env.get('CLIENT_SECRET')}&grant_type=client_credentials`, {
    method: 'POST',
  });
  const IGDBAuthData = await IGDBAuthResponse.json();

  const body = `search "${name}"; limit 100; fields name, cover.image_id; where version_parent = null; where parent_game = null;`

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
    { status: 200, headers: { "Content-Type": "application/json" } },
  )
})

/* To invoke locally:

  1. Run `supabase start` (see: https://supabase.com/docs/reference/cli/supabase-start)
  2. Make an HTTP request:

  curl -i --location --request POST 'http://127.0.0.1:54321/functions/v1/search-games' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0' \
    --header 'Content-Type: application/json' \
    --data '{"name":"Functions"}'

*/
