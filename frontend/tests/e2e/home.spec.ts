import { test, expect } from '@playwright/test'

test('home page loads and shows navigation', async ({ page }) => {
  await page.goto('/')
  await expect(page).toHaveTitle(/MyAnimeAPI/i)
  await expect(page.getByRole('navigation')).toBeVisible()
})

test('animes page is reachable', async ({ page }) => {
  await page.goto('/animes')
  await expect(page.getByRole('heading')).toBeVisible()
})
