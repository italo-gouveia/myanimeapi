import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AdminAnalyticsPage } from './AdminAnalyticsPage'

vi.mock('@tanstack/react-query', () => ({ useQuery: vi.fn() }))
vi.mock('../../lib/api/export', () => ({ getAdminAnalytics: vi.fn() }))

import { useQuery } from '@tanstack/react-query'

const mockAnalytics = {
  recent_reviews_per_day: [
    { date: '2024-11-10', count: 5 },
    { date: '2024-11-11', count: 8 },
    { date: '2024-11-12', count: 3 },
    { date: '2024-11-13', count: 12 },
    { date: '2024-11-14', count: 7 },
    { date: '2024-11-15', count: 9 },
    { date: '2024-11-16', count: 4 },
  ],
  recent_signups_per_day: [
    { date: '2024-11-10', count: 2 },
    { date: '2024-11-11', count: 1 },
    { date: '2024-11-12', count: 3 },
    { date: '2024-11-13', count: 0 },
    { date: '2024-11-14', count: 5 },
    { date: '2024-11-15', count: 2 },
    { date: '2024-11-16', count: 1 },
  ],
  top_favorited: [
    { anime_id: 1, title: 'Naruto', count: 150, avg_rating: null },
    { anime_id: 2, title: 'One Piece', count: 120, avg_rating: null },
  ],
  top_reviewed: [
    { anime_id: 3, title: 'Bleach', count: 80, avg_rating: 8.5 },
    { anime_id: 4, title: 'FMA', count: 60, avg_rating: 9.1 },
  ],
}

function renderPage() {
  return render(<MemoryRouter><AdminAnalyticsPage /></MemoryRouter>)
}

describe('AdminAnalyticsPage', () => {
  it('renders Analytics heading', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: true, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('heading', { name: /analytics/i })).toBeInTheDocument()
  })

  it('renders SparkBar labels while loading', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: true, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/reviews — last 7 days/i)).toBeInTheDocument()
    expect(screen.getByText(/signups — last 7 days/i)).toBeInTheDocument()
  })

  it('renders TopList headings while loading', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: true, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/most favorited/i)).toBeInTheDocument()
    expect(screen.getByText(/most reviewed/i)).toBeInTheDocument()
  })

  it('renders top favorited anime links when data is loaded', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: mockAnalytics, isLoading: false, isError: false,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('link', { name: /naruto/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /one piece/i })).toBeInTheDocument()
  })

  it('renders top reviewed anime with ratings', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: mockAnalytics, isLoading: false, isError: false,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('link', { name: /bleach/i })).toBeInTheDocument()
    expect(screen.getByText(/8.5/)).toBeInTheDocument()
    expect(screen.getByText(/9.1/)).toBeInTheDocument()
  })

  it('shows error alert on failure', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: false, isError: true } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/failed to load analytics/i)).toBeInTheDocument()
  })

  it('renders "No data yet" for empty lists', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: {
        recent_reviews_per_day: [],
        recent_signups_per_day: [],
        top_favorited: [],
        top_reviewed: [],
      },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getAllByText(/no data yet/i).length).toBeGreaterThanOrEqual(2)
  })

  it('renders bar counts when data is loaded', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: mockAnalytics, isLoading: false, isError: false,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('150')).toBeInTheDocument()
    expect(screen.getByText('120')).toBeInTheDocument()
    expect(screen.getByText('80')).toBeInTheDocument()
  })
})
